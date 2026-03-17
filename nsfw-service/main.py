from typing import Any, Dict

import io
import json
import os
import time
from http.client import HTTPConnection
from urllib.parse import urlparse

from fastapi import FastAPI, File, HTTPException, UploadFile
from fastapi.responses import JSONResponse
from PIL import Image, UnidentifiedImageError
from transformers import Pipeline, pipeline


app = FastAPI(title="NSFW Detection Service")

nsfw_classifier: Pipeline | None = None
nsfw_threshold = float(os.getenv("NSFW_THRESHOLD", "0.7"))
model_id = os.getenv("NSFW_MODEL_ID", "Falconsai/nsfw_image_detection")


def _write_debug_log(hypothesis_id: str, location: str, message: str, data: Dict[str, Any] | None = None) -> None:
    # region agent log
    entry: Dict[str, Any] = {
        "sessionId": "0f9762",
        "runId": "initial",
        "hypothesisId": hypothesis_id,
        "location": location,
        "message": message,
        "data": data or {},
        "timestamp": int(time.time() * 1000),
    }

    body = json.dumps(entry)

    try:
        url = urlparse("http://host.docker.internal:7586/ingest/e28d96d1-7f10-4f83-951e-74e65210567a")
        conn = HTTPConnection(url.hostname, url.port or 80, timeout=1.0)
        conn.request(
            "POST",
            url.path,
            body=body,
            headers={
                "Content-Type": "application/json",
                "X-Debug-Session-Id": "0f9762",
            },
        )
        conn.getresponse().read()
        conn.close()
    except OSError:
        # Logging must never break the service
        return
    # endregion


@app.on_event("startup")
def load_model() -> None:
    global nsfw_classifier

    try:
        nsfw_classifier = pipeline("image-classification", model=model_id)
    except Exception as exc:  # noqa: BLE001
        nsfw_classifier = None
        raise RuntimeError(f"Failed to load NSFW model '{model_id}': {exc}") from exc


@app.post("/classify")
async def classify_image(file: UploadFile = File(...)) -> JSONResponse:
    _write_debug_log(
        "C",
        "main.py:classify",
        "received_image_for_classification",
        {"filename": file.filename, "content_type": file.content_type},
    )

    if nsfw_classifier is None:
        raise HTTPException(status_code=503, detail="NSFW classifier not initialized")

    # Accept standard image/* types and the default application/octet-stream
    # that the backend may send when proxying the upload.
    if (
        file.content_type
        and not file.content_type.startswith("image/")
        and file.content_type != "application/octet-stream"
    ):
        raise HTTPException(status_code=400, detail="Only image uploads are supported")

    try:
        data = await file.read()

        if not data:
            raise HTTPException(status_code=400, detail="Empty file")

        image = Image.open(io.BytesIO(data)).convert("RGB")
    except UnidentifiedImageError as exc:
        _write_debug_log(
            "C",
            "main.py:classify",
            "invalid_image_data",
            {"filename": file.filename},
        )
        raise HTTPException(status_code=400, detail="Invalid image data") from exc

    try:
        results: list[Dict[str, Any]] = nsfw_classifier(image)  # type: ignore[misc]
    except Exception as exc:  # noqa: BLE001
        _write_debug_log(
            "D",
            "main.py:classify",
            "nsfw_classification_failed",
            {"error": str(exc)},
        )
        raise HTTPException(status_code=500, detail="Failed to run NSFW classification") from exc

    if not results:
        _write_debug_log(
            "D",
            "main.py:classify",
            "no_classification_results",
            {"filename": file.filename},
        )
        raise HTTPException(status_code=500, detail="No classification results returned")

    top = max(results, key=lambda r: float(r.get("score", 0.0)))
    label = str(top.get("label", "unknown"))
    score = float(top.get("score", 0.0))

    # Consider an image NSFW only when the predicted label explicitly
    # indicates NSFW content and its score passes the threshold.
    is_nsfw = "nsfw" in label.lower() and score >= nsfw_threshold

    payload = {
        "is_nsfw": is_nsfw,
        "label": label,
        "score": score,
    }

    _write_debug_log(
        "D",
        "main.py:classify",
        "nsfw_classification_result",
        {
            "filename": file.filename,
            "label": label,
            "score": score,
            "is_nsfw": is_nsfw,
            "threshold": nsfw_threshold,
        },
    )

    return JSONResponse(content=payload)


@app.get("/health")
def health() -> Dict[str, Any]:
    return {
        "status": "ok",
        "model_loaded": nsfw_classifier is not None,
        "model_id": model_id,
        "threshold": nsfw_threshold,
    }

