import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { BrowserRouter } from 'react-router-dom'
import Landing from './Landing'

describe('Landing Page', () => {
  it('renders the main heading', () => {
    render(
      <BrowserRouter>
        <Landing />
      </BrowserRouter>
    )
    
    expect(screen.getByText(/Find Your/i)).toBeInTheDocument()
    expect(screen.getByText(/Perfect Match/i)).toBeInTheDocument()
  })

  it('renders Create Account button', () => {
    render(
      <BrowserRouter>
        <Landing />
      </BrowserRouter>
    )
    
    expect(screen.getByText('Create Account')).toBeInTheDocument()
  })

  it('renders Sign In button', () => {
    render(
      <BrowserRouter>
        <Landing />
      </BrowserRouter>
    )
    
    expect(screen.getByText('Sign In')).toBeInTheDocument()
  })

  it('renders feature cards', () => {
    render(
      <BrowserRouter>
        <Landing />
      </BrowserRouter>
    )
    
    expect(screen.getByText('Match')).toBeInTheDocument()
    expect(screen.getByText('Chat')).toBeInTheDocument()
    expect(screen.getByText('Connect')).toBeInTheDocument()
  })
})
