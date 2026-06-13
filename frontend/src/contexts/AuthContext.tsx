'use client'

import { createContext, useContext, useState, useEffect, useCallback, type ReactNode } from 'react'
import { getToken, setSession, clearSession, isAuthenticated as checkAuth } from '@/lib/auth'
import { loginUser, registerUser, type AuthResponse } from '@/lib/api'

interface UserInfo {
  id: string
  type: string
}

interface AuthContextType {
  token: string | null
  user: UserInfo | null
  isAuthenticated: boolean
  login: (email: string, password: string) => Promise<void>
  register: (name: string, email: string, password: string) => Promise<void>
  logout: () => void
}

const AuthContext = createContext<AuthContextType | null>(null)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [token, setToken] = useState<string | null>(null)
  const [user, setUser] = useState<UserInfo | null>(null)

  useEffect(() => {
    if (checkAuth()) {
      setToken(getToken())
      try {
        const stored = localStorage.getItem('harness_user')
        if (stored) setUser(JSON.parse(stored))
      } catch {}
    }
  }, [])

  const login = useCallback(async (email: string, password: string) => {
    const resp: AuthResponse = await loginUser(email, password)
    setSession(resp.token, resp.user_id, resp.user_type)
    setToken(resp.token)
    setUser({ id: resp.user_id, type: resp.user_type })
  }, [])

  const register = useCallback(async (name: string, email: string, password: string) => {
    await registerUser(name, email, password)
  }, [])

  const logout = useCallback(() => {
    clearSession()
    setToken(null)
    setUser(null)
  }, [])

  return (
    <AuthContext.Provider value={{ token, user, isAuthenticated: !!token, login, register, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth(): AuthContextType {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}
