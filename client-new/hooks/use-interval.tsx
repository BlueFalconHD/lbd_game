import { useEffect, useRef } from "react";

const useInterval = (fn: () => void, delay: number | null) => {
  const fnRef = useRef<() => void>(null);

  useEffect(() => {
    fnRef.current = fn;
  }, [fn]);

  useEffect(() => {
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
