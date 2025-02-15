export const getCookie = (name: string): string | null => {
  if (typeof window === "undefined") return null;
  const cookie = document.cookie
    .split("; ")
    .find((row) => row.startsWith(name));
  return cookie ? cookie.split("=")[1] : null;
};

export const deleteCookie = (name: string): void => {
  if (typeof window === "undefined") return;
  document.cookie = name + "=; Max-Age=-99999999;";
};
