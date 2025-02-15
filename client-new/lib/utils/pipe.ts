export const pipe = <T>(value: T, label?: string): T => {
  if (label) {
    console.log(label, value);
  } else {
    console.log(value);
  }
  return value;
};
