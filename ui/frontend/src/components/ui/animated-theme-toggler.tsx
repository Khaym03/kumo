'use client'

import { Moon, SunDim } from 'lucide-react'
import { useState, useRef } from 'react'
import { flushSync } from 'react-dom'
import { cn } from '@/lib/utils'
import { useAppContext } from '@/context/app_ctx'
import { Button } from './button'

type props = {
  className?: string
}

export const AnimatedThemeToggler = ({ className }: props) => {
  const { isDarkMode, setIsDarkMode } = useAppContext()

  const buttonRef = useRef<HTMLButtonElement | null>(null)
  const changeTheme = async () => {
    if (!buttonRef.current) return

    await document.startViewTransition(() => {
      flushSync(() => {
        const dark = document.documentElement.classList.toggle('dark')
        setIsDarkMode(dark)
      })
    }).ready

    const { top, left, width, height } =
      buttonRef.current.getBoundingClientRect()
    const y = top + height / 2
    const x = left + width / 2

    const right = window.innerWidth - left
    const bottom = window.innerHeight - top
    const maxRad = Math.hypot(Math.max(left, right), Math.max(top, bottom))

    document.documentElement.animate(
      {
        clipPath: [
          `circle(0px at ${x}px ${y}px)`,
          `circle(${maxRad}px at ${x}px ${y}px)`
        ]
      },
      {
        duration: 700,
        easing: 'ease-in-out',
        pseudoElement: '::view-transition-new(root)'
      }
    )
  }
  return (
    <Button
      ref={buttonRef}
      variant="ghost"
      size="icon"
      onClick={changeTheme}
      className={(cn(className), 'size-8 rounded-none')}
    >
      {isDarkMode ? <SunDim /> : <Moon />}
    </Button>
  )
}
