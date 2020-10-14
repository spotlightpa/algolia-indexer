export default function formatTel(s) {
  if (!s) {
    return "";
  }
  let re = /^tel:1?(\d{3})(\d{3})(\d{4})$/.exec(s);
  if (re) {
    return `(${re[1]}) ${re[2]}-${re[3]}`;
  }
  return s;
}
