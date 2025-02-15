import { useEffect, useRef } from "react";

const useTimeout = (fn: () => void, delay: number | null) => {
  const fnRef = useRef<() => void>(null);

  useEffect(() => {
    fnRef.current = fn;
  }, [fn]);

  useEffect(() => {
    if (delay === null) return;

    const handle = setTimeout(() => {
      if (fnRef.current) {
        fnRef.current();
      }
    }, delay);

    return () => clearTimeout(handle);
  }, [delay]);

  return;
};

export default useTimeout;
