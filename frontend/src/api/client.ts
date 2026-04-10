import { factory } from './api'

const RPC_URL = '/v1/rpc/'

class RpcError extends Error {
  constructor(
    public code: number,
    message: string,
  ) {
    super(message)
    this.name = 'RpcError'
  }
}

let _onUnauthorized: (() => void) | null = null

export function setUnauthorizedHandler(fn: () => void) {
  _onUnauthorized = fn
}

async function send(method: string, params?: unknown): Promise<unknown> {
  const token = localStorage.getItem('token')

  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  }
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  const response = await fetch(RPC_URL, {
    method: 'POST',
    headers,
    body: JSON.stringify({
      jsonrpc: '2.0',
      id: 1,
      method,
      params: params ?? {},
    }),
  })

  if (response.status === 401) {
    localStorage.removeItem('token')
    _onUnauthorized?.()
    throw new RpcError(401, 'Unauthorized')
  }

  const json = await response.json()

  if (json.error) {
    throw new RpcError(json.error.code, json.error.message)
  }

  return json.result
}

export const api = factory(send)
export { RpcError }
