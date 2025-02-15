import { useEffect, useRef } from "react";

const useInterval = (fn: () => void, delay: number | undefined) => {
  const fnRef = useRef<() => void>(fn);

  useEffect(() => {
    fnRef.current = fn;
  }, [fn]);

  useEffect(() => {
    if (delay === undefined) return;

    const interval = setInterval(() => {
      if (fnRef.current) {
        fnRef.current();
      }
    }, delay);

    return () => clearInterval(interval);
  }, [delay]);

  return;
};

export default useInterval;
