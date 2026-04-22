import crypto from "node:crypto";

function extractDomain(url) {
  try {
    return new URL(url).hostname.replace(/^www\./, "");
  } catch {
    return (url || "").replace(/^www\./, "");
  }
}

function generateSec(projectName) {
  const key = "04788f1ea15d305f"; // 必须与 fuck.js 里的 key 一致
  const domain = extractDomain(projectName);
  const projectNameMd5 = crypto.createHash("md5").update(domain).digest("hex");

  const maxLen = Math.max(projectNameMd5.length, key.length);
  let sec = "";

  for (let i = 0; i < maxLen; i++) {
    sec += projectNameMd5[i] || "x"; // i%3===0
    sec += key[i] || "x";            // i%3===1
    sec += "x";                      // i%3===2 filler
  }

  return { domain, projectNameMd5, key, sec };
}

console.log(generateSec("https://icosmos.space"));