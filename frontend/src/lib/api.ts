const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'

interface RequestOptions {
  method?: string
  body?: unknown
  token?: string
}

export async function apiRequest<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  }

  if (options.token) {
    headers['Authorization'] = `Bearer ${options.token}`
  }

  const response = await fetch(`${API_BASE}${path}`, {
    method: options.method || 'GET',
    headers,
    body: options.body ? JSON.stringify(options.body) : undefined,
  })

  if (!response.ok) {
    const errorBody = await response.json().catch(() => ({ error: 'Unknown error' }))
    throw new Error(errorBody.error || `HTTP ${response.status}`)
  }

  return response.json()
}

export interface AuthResponse {
  token: string
  user_id: string
  user_type: 'human' | 'ai'
  expires_at: number
}

export interface UserResponse {
  id: string
  name: string
  email: string
  avatar_url?: string
  created_at: string
  updated_at: string
}

export interface AIAgent {
  id: string
  name: string
  model_type: string
  capabilities: string[]
  permission_level: string
  metadata: Record<string, unknown>
  is_active: boolean
  created_at: string
  updated_at: string
}

export async function loginUser(email: string, password: string): Promise<AuthResponse> {
  return apiRequest<AuthResponse>('/auth/login', {
    method: 'POST',
    body: { email, password },
  })
}

export async function registerUser(name: string, email: string, password: string): Promise<UserResponse> {
  return apiRequest<UserResponse>('/auth/register', {
    method: 'POST',
    body: { name, email, password },
  })
}

export async function registerAgent(data: {
  name: string
  model_type: string
  capabilities?: string[]
  permission_level?: string
  metadata?: Record<string, unknown>
}): Promise<{ agent: AIAgent; api_key: string }> {
  return apiRequest('/auth/agents/register', {
    method: 'POST',
    body: data,
  })
}

export async function authenticateAgent(agentId: string, apiKey: string): Promise<AuthResponse> {
  return apiRequest<AuthResponse>('/auth/agents/auth', {
    method: 'POST',
    body: { agent_id: agentId, api_key: apiKey },
  })
}
