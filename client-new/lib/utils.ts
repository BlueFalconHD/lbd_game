import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function unixTime5SecondsFromNow() {
  return Math.floor(Date.now() / 1000) + 5;
}
