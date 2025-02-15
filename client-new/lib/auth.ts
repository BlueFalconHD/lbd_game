export function GetToken(): string | null {
  return localStorage.getItem("token");
}

export function SetToken(token: string) {
  localStorage.setItem("token", token);
}

export function DeleteToken() {
  localStorage.removeItem("token");
}
