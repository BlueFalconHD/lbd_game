import { useEffect, useRef } from "react";

const useTimeout = (fn: () => void, delay: number | undefined) => {
  const fnRef = useRef<() => void>(fn);

  useEffect(() => {
    fnRef.current = fn;
  }, [fn]);

  useEffect(() => {
    if (delay === undefined) return;

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
