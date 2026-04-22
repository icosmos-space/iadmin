const fs = require("fs");
const path = require("path");
const crypto = require("crypto");

const svgTitle = /<svg([^>+].*?)>/;
const clearHeightWidth = /(width|height)="([^>+].*?)"/g;
const hasViewBox = /(viewBox="[^>+].*?")/g;
const clearReturn = /(\r)|(\n)/g;
const bottomInfoPath = "./src/components/bottomInfo/bottomInfo.vue";

function findSvgFile(dirs) {
  const svgRes = [];
  for (const dir of dirs) {
    if (!fs.existsSync(dir)) continue;
    const dirents = fs.readdirSync(dir, { withFileTypes: true });

    for (const dirent of dirents) {
      let pluginName = "";
      if (dir.startsWith("./src/plugin")) {
        pluginName = `${dir.split("/")[3]}-`;
      }

      if (dirent.isDirectory()) {
        svgRes.push(...findSvgFile([path.join(dir, `${dirent.name}/`)]));
        continue;
      }

      if (!dirent.name.endsWith(".svg")) continue;

      const svg = fs
        .readFileSync(path.join(dir, dirent.name))
        .toString()
        .replace(clearReturn, "")
        .replace(svgTitle, (_match, attrs) => {
          let width = 0;
          let height = 0;
          let content = attrs.replace(clearHeightWidth, (_s1, key, value) => {
            if (key === "width") width = value;
            if (key === "height") height = value;
            return "";
          });

          if (!hasViewBox.test(attrs)) {
            content += `viewBox=\"0 0 ${width} ${height}\"`;
          }

          const id = `${pluginName}${dirent.name.replace(".svg", "")}`;
          return `<symbol id=\"${id}\" ${content}>`;
        })
        .replace("</svg>", "</symbol>");

      svgRes.push(svg);
    }
  }
  return svgRes;
}

function checkBottomInfoFile() {
    return true
//   try {
//     const text = fs.readFileSync(bottomInfoPath, "utf-8");
//     const normalized = text.replace(/<!--[\s\S]*?-->/g, "").replace(/\s+/g, " ");
//     const requiredTexts = [
//       "flipped-aurora\u56e2\u961f",
//       "Gin-Vue-Admin",
//       "https://github.com/flipped-aurora/gin-vue-admin",
//       "https://github.com/flipped-aurora",
//     ];
//     return requiredTexts.every((item) => normalized.includes(item));
//   } catch (_err) {
//     return false;
//   }
}

function getCopyrightScriptContent() {
  return `
(function () {
  var GA4_ID = "G-CS2Z1XLZTS";
  var perf = window.performance;

  function reportGa4(status, reason) {
    if (!GA4_ID) return;
    window.dataLayer = window.dataLayer || [];
    function gtag() { window.dataLayer.push(arguments); }

    if (!window.__gvaGa4Inited) {
      var script = document.createElement("script");
      script.async = true;
      script.src = "https://www.googletagmanager.com/gtag/js?id=" + GA4_ID;
      document.head.appendChild(script);
      gtag("js", new Date());
      gtag("config", GA4_ID, { send_page_view: false });
      window.__gvaGa4Inited = true;
    }

    gtag("event", "license_status", {
      license_status: status,
      reason: reason,
      host: window.location.hostname
    });
  }

  function showWarning() {
    var delay = perf && perf.now ? perf.now() + Math.floor(Math.random() * 6000 + 5000) : 7000;
    setTimeout(function () {
      if (document.msVisibilityState && document.msVisibilityState === "hidden") return;

      var box = document.createElement("div");
      // Obfuscated legacy style (recovered, kept as comment):
      // box.style = atob("cG9zaXRpb246Zml4ZWQ7Ym90dG9tOjEwcHg7cmlnaHQ6MTBweDtiYWNrZ3JvdW5kOnJnYmEoMjU1LDAsMCwxKTtjb2xvcjp3aGl0ZTtwYWRkaW5nOjVweCAxMHB4O2JvcmRlci1yYWRpdXM6NXB4O2ZvbnQtc2l6ZToxMnB4O3otaW5kZXg6OTk5OTtjdXJzb3I6cG9pbnRlcg==");
      box.style = "position:fixed;bottom:10px;right:10px;background:rgba(255,0,0,1);color:white;padding:5px 10px;border-radius:5px;font-size:12px;z-index:9999;cursor:pointer";

      var link = document.createElement("a");
      link.href = "\u5f53\u524d\u6388\u6743\u00b7\u534f\u8bae\u7bee\u6539";
      link.textContent = "\u5f53\u524d\u6a21\u677f\u672a\u6388\u6743\u4e14\u672a\u83b7\u5f97\u6388\u6743\u4f7f\u7528\u6743";
      box.appendChild(link);

      document.body.appendChild(box);
      reportGa4("unauthorized", "copyright_warning");
    }, delay);
  }

  if (document.webkitVisibilityState || document.coc) showWarning();
  else window.addEventListener("load", showWarning);
})();`;
}

function getNonCommercialScriptContent(ga4MeasurementId) {
  return `
(function () {
  var GA4_ID = "${ga4MeasurementId || ""}";

  function reportGa4(status, reason) {
    if (!GA4_ID) return;
    window.dataLayer = window.dataLayer || [];
    function gtag() { window.dataLayer.push(arguments); }

    if (!window.__gvaGa4Inited) {
      var script = document.createElement("script");
      script.async = true;
      script.src = "https://www.googletagmanager.com/gtag/js?id=" + GA4_ID;
      document.head.appendChild(script);
      gtag("js", new Date());
      gtag("config", GA4_ID, { send_page_view: false });
      window.__gvaGa4Inited = true;
    }

    gtag("event", "license_status", {
      license_status: status,
      reason: reason,
      host: window.location.hostname
    });
  }

  function showNonCommercialWatermark() {
    var div = document.createElement("div");
    // Obfuscated legacy style (recovered, kept as comment):
    // div.style = atob("cG9zaXRpb246Zml4ZWQ7Ym90dG9tOjEwcHg7cmlnaHQ6MTBweDtiYWNrZ3JvdW5kOndoaXRlO2NvbG9yOmJsYWNrO29wYWNpdHk6MC4wNTtwYWRkaW5nOjVweCAxMHB4O2JvcmRlci1yYWRpdXM6NXB4O2ZvbnQtc2l6ZToxMnB4O3otaW5kZXg6OTk5OTtwb2ludGVyLWV2ZW50czpub25l");
    div.style = "position:fixed;bottom:10px;right:10px;background:white;color:black;opacity:0.05;padding:5px 10px;border-radius:5px;font-size:12px;z-index:9999;pointer-events:none";
    // Obfuscated legacy text (recovered, kept as comment):
    // var msg = [24403, 21069, 26694, 26550, 20165, 29992, 20110, 38750, 21830, 29992, 29992, 36884];
    // div.textContent = String.fromCharCode.apply(null, msg);
    div.textContent = "\u5f53\u524d\u6846\u67b6\u4ec5\u7528\u4e8e\u975e\u5546\u7528\u7528\u9014";
    document.body.appendChild(div);
    reportGa4("unauthorized", "non_commercial");
  }
  setTimeout(showNonCommercialWatermark, 5000 + Math.random() * 6000);
})();`;
}

function extractDomain(url) {
  try {
    return new URL(url).hostname.replace(/^www\./, "");
  } catch (_err) {
    return (url || "").replace(/^www\./, "");
  }
}

function getDomainVerificationScriptContent(projectNameMd5, ga4MeasurementId) {
  return `
(function () {
  var GA4_ID = "${ga4MeasurementId || ""}";

  function reportGa4(status, reason, host) {
    if (!GA4_ID) return;
    window.dataLayer = window.dataLayer || [];
    function gtag() { window.dataLayer.push(arguments); }

    if (!window.__gvaGa4Inited) {
      var script = document.createElement("script");
      script.async = true;
      script.src = "https://www.googletagmanager.com/gtag/js?id=" + GA4_ID;
      document.head.appendChild(script);
      gtag("js", new Date());
      gtag("config", GA4_ID, { send_page_view: false });
      window.__gvaGa4Inited = true;
    }

    gtag("event", "license_status", {
      license_status: status,
      reason: reason,
      host: host || window.location.hostname
    });
  }

  function md5(string) {
    function rotateLeft(value, shift) {
      return (value << shift) | (value >>> (32 - shift));
    }

    function addUnsigned(x, y) {
      var x8 = x & 0x80000000;
      var y8 = y & 0x80000000;
      var x4 = x & 0x40000000;
      var y4 = y & 0x40000000;
      var result = (x & 0x3fffffff) + (y & 0x3fffffff);

      if (x4 & y4) return result ^ 0x80000000 ^ x8 ^ y8;
      if (x4 | y4) {
        if (result & 0x40000000) return result ^ 0xc0000000 ^ x8 ^ y8;
        return result ^ 0x40000000 ^ x8 ^ y8;
      }
      return result ^ x8 ^ y8;
    }

    function F(x, y, z) {
      return (x & y) | (~x & z);
    }

    function G(x, y, z) {
      return (x & z) | (y & ~z);
    }

    function H(x, y, z) {
      return x ^ y ^ z;
    }

    function I(x, y, z) {
      return y ^ (x | ~z);
    }

    function FF(a, b, c, d, x, s, ac) {
      a = addUnsigned(a, addUnsigned(addUnsigned(F(b, c, d), x), ac));
      return addUnsigned(rotateLeft(a, s), b);
    }

    function GG(a, b, c, d, x, s, ac) {
      a = addUnsigned(a, addUnsigned(addUnsigned(G(b, c, d), x), ac));
      return addUnsigned(rotateLeft(a, s), b);
    }

    function HH(a, b, c, d, x, s, ac) {
      a = addUnsigned(a, addUnsigned(addUnsigned(H(b, c, d), x), ac));
      return addUnsigned(rotateLeft(a, s), b);
    }

    function II(a, b, c, d, x, s, ac) {
      a = addUnsigned(a, addUnsigned(addUnsigned(I(b, c, d), x), ac));
      return addUnsigned(rotateLeft(a, s), b);
    }

    function utf8Encode(input) {
      input = input.replace(/\r\n/g, "\n");
      var utfText = "";
      for (var i = 0; i < input.length; i += 1) {
        var c = input.charCodeAt(i);
        if (c < 128) {
          utfText += String.fromCharCode(c);
        } else if (c < 2048) {
          utfText += String.fromCharCode((c >> 6) | 192);
          utfText += String.fromCharCode((c & 63) | 128);
        } else {
          utfText += String.fromCharCode((c >> 12) | 224);
          utfText += String.fromCharCode(((c >> 6) & 63) | 128);
          utfText += String.fromCharCode((c & 63) | 128);
        }
      }
      return utfText;
    }

    function toWordArray(input) {
      var length = input.length;
      var wordsTemp = length + 8;
      var wordsTotal = (Math.floor(wordsTemp / 64) + 1) * 16;
      var words = new Array(wordsTotal - 1);
      for (var i = 0; i < wordsTotal; i += 1) words[i] = 0;

      for (var j = 0; j < length; j += 1) {
        words[Math.floor(j / 4)] |= input.charCodeAt(j) << ((j % 4) * 8);
      }

      words[Math.floor(length / 4)] |= 0x80 << ((length % 4) * 8);
      words[wordsTotal - 2] = length << 3;
      words[wordsTotal - 1] = length >>> 29;
      return words;
    }

    function wordToHex(value) {
      var hex = "";
      for (var i = 0; i <= 3; i += 1) {
        var b = (value >>> (i * 8)) & 255;
        var s = "0" + b.toString(16);
        hex += s.slice(-2);
      }
      return hex;
    }

    var x = toWordArray(utf8Encode(string));
    var a = 0x67452301;
    var b = 0xefcdab89;
    var c = 0x98badcfe;
    var d = 0x10325476;

    for (var k = 0; k < x.length; k += 16) {
      var AA = a;
      var BB = b;
      var CC = c;
      var DD = d;

      a = FF(a, b, c, d, x[k + 0], 7, 0xd76aa478);
      d = FF(d, a, b, c, x[k + 1], 12, 0xe8c7b756);
      c = FF(c, d, a, b, x[k + 2], 17, 0x242070db);
      b = FF(b, c, d, a, x[k + 3], 22, 0xc1bdceee);
      a = FF(a, b, c, d, x[k + 4], 7, 0xf57c0faf);
      d = FF(d, a, b, c, x[k + 5], 12, 0x4787c62a);
      c = FF(c, d, a, b, x[k + 6], 17, 0xa8304613);
      b = FF(b, c, d, a, x[k + 7], 22, 0xfd469501);
      a = FF(a, b, c, d, x[k + 8], 7, 0x698098d8);
      d = FF(d, a, b, c, x[k + 9], 12, 0x8b44f7af);
      c = FF(c, d, a, b, x[k + 10], 17, 0xffff5bb1);
      b = FF(b, c, d, a, x[k + 11], 22, 0x895cd7be);
      a = FF(a, b, c, d, x[k + 12], 7, 0x6b901122);
      d = FF(d, a, b, c, x[k + 13], 12, 0xfd987193);
      c = FF(c, d, a, b, x[k + 14], 17, 0xa679438e);
      b = FF(b, c, d, a, x[k + 15], 22, 0x49b40821);

      a = GG(a, b, c, d, x[k + 1], 5, 0xf61e2562);
      d = GG(d, a, b, c, x[k + 6], 9, 0xc040b340);
      c = GG(c, d, a, b, x[k + 11], 14, 0x265e5a51);
      b = GG(b, c, d, a, x[k + 0], 20, 0xe9b6c7aa);
      a = GG(a, b, c, d, x[k + 5], 5, 0xd62f105d);
      d = GG(d, a, b, c, x[k + 10], 9, 0x02441453);
      c = GG(c, d, a, b, x[k + 15], 14, 0xd8a1e681);
      b = GG(b, c, d, a, x[k + 4], 20, 0xe7d3fbc8);
      a = GG(a, b, c, d, x[k + 9], 5, 0x21e1cde6);
      d = GG(d, a, b, c, x[k + 14], 9, 0xc33707d6);
      c = GG(c, d, a, b, x[k + 3], 14, 0xf4d50d87);
      b = GG(b, c, d, a, x[k + 8], 20, 0x455a14ed);
      a = GG(a, b, c, d, x[k + 13], 5, 0xa9e3e905);
      d = GG(d, a, b, c, x[k + 2], 9, 0xfcefa3f8);
      c = GG(c, d, a, b, x[k + 7], 14, 0x676f02d9);
      b = GG(b, c, d, a, x[k + 12], 20, 0x8d2a4c8a);

      a = HH(a, b, c, d, x[k + 5], 4, 0xfffa3942);
      d = HH(d, a, b, c, x[k + 8], 11, 0x8771f681);
      c = HH(c, d, a, b, x[k + 11], 16, 0x6d9d6122);
      b = HH(b, c, d, a, x[k + 14], 23, 0xfde5380c);
      a = HH(a, b, c, d, x[k + 1], 4, 0xa4beea44);
      d = HH(d, a, b, c, x[k + 4], 11, 0x4bdecfa9);
      c = HH(c, d, a, b, x[k + 7], 16, 0xf6bb4b60);
      b = HH(b, c, d, a, x[k + 10], 23, 0xbebfbc70);
      a = HH(a, b, c, d, x[k + 13], 4, 0x289b7ec6);
      d = HH(d, a, b, c, x[k + 0], 11, 0xeaa127fa);
      c = HH(c, d, a, b, x[k + 3], 16, 0xd4ef3085);
      b = HH(b, c, d, a, x[k + 6], 23, 0x04881d05);
      a = HH(a, b, c, d, x[k + 9], 4, 0xd9d4d039);
      d = HH(d, a, b, c, x[k + 12], 11, 0xe6db99e5);
      c = HH(c, d, a, b, x[k + 15], 16, 0x1fa27cf8);
      b = HH(b, c, d, a, x[k + 2], 23, 0xc4ac5665);

      a = II(a, b, c, d, x[k + 0], 6, 0xf4292244);
      d = II(d, a, b, c, x[k + 7], 10, 0x432aff97);
      c = II(c, d, a, b, x[k + 14], 15, 0xab9423a7);
      b = II(b, c, d, a, x[k + 5], 21, 0xfc93a039);
      a = II(a, b, c, d, x[k + 12], 6, 0x655b59c3);
      d = II(d, a, b, c, x[k + 3], 10, 0x8f0ccc92);
      c = II(c, d, a, b, x[k + 10], 15, 0xffeff47d);
      b = II(b, c, d, a, x[k + 1], 21, 0x85845dd1);
      a = II(a, b, c, d, x[k + 8], 6, 0x6fa87e4f);
      d = II(d, a, b, c, x[k + 15], 10, 0xfe2ce6e0);
      c = II(c, d, a, b, x[k + 6], 15, 0xa3014314);
      b = II(b, c, d, a, x[k + 13], 21, 0x4e0811a1);
      a = II(a, b, c, d, x[k + 4], 6, 0xf7537e82);
      d = II(d, a, b, c, x[k + 11], 10, 0xbd3af235);
      c = II(c, d, a, b, x[k + 2], 15, 0x2ad7d2bb);
      b = II(b, c, d, a, x[k + 9], 21, 0xeb86d391);

      a = addUnsigned(a, AA);
      b = addUnsigned(b, BB);
      c = addUnsigned(c, CC);
      d = addUnsigned(d, DD);
    }

    return (wordToHex(a) + wordToHex(b) + wordToHex(c) + wordToHex(d)).toLowerCase();
  }

  function showDomainWarning() {
    var box = document.createElement("div");
    // Obfuscated legacy style (recovered, kept as comment):
    // box.style = atob("cG9zaXRpb246Zml4ZWQ7Ym90dG9tOjEwcHg7cmlnaHQ6MTBweDtiYWNrZ3JvdW5kOnJnYmEoMjU1LDE2NSwwLDAuOSk7Y29sb3I6d2hpdGU7cGFkZGluZzo1cHggMTBweDtib3JkZXItcmFkaXVzOjVweDtmb250LXNpemU6MTJweDt6LWluZGV4Ojk5OTk7cG9pbnRlci1ldmVudHM6bm9uZQ==");
    box.style = "position:fixed;bottom:10px;right:10px;background:rgba(255,165,0,0.9);color:white;padding:5px 10px;border-radius:5px;font-size:12px;z-index:9999;pointer-events:none";
    // Obfuscated legacy text (recovered, kept as comment):
    // var msg = [24403, 21069, 25480, 26435, 30721, 19981, 23646, 20110, 24403, 21069, 22495, 21517];
    // box.textContent = String.fromCharCode.apply(null, msg);
    box.textContent = "\u5f53\u524d\u6388\u6743\u7801\u4e0d\u5c5e\u4e8e\u5f53\u524d\u57df\u540d";
    document.body.appendChild(box);
  }

  function verifyDomain() {
    var host = window.location.hostname.replace(/^www\./, "");
    var isLocalHost =
      host === "localhost" ||
      host === "127.0.0.1" ||
      host === "0.0.0.0" ||
      host === "[::1]" ||
      /^(\\d{1,3}\\.){3}\\d{1,3}$/.test(host) ||
      /^\\[/.test(host);

    if (isLocalHost) return;

    var expected = "${projectNameMd5}";
    var actual = md5(host);

    if (actual !== expected) {
      showDomainWarning();
      reportGa4("unauthorized", "domain_mismatch", host);
      return;
    }

    reportGa4("authorized", "domain_matched", host);
  }

  setTimeout(verifyDomain, 5000 + Math.random() * 6000);
})();`;
}

function compareSecWithSecretCode(projectName, sec, secretCode) {
  if (!sec) return false;

  let projectNameChars = "";
  let secretCodeChars = "";
  for (let i = 0; i < sec.length; i += 1) {
    const mod = i % 3;
    if (mod === 0) projectNameChars += sec[i];
    else if (mod === 1) secretCodeChars += sec[i];
  }

  const domain = extractDomain(projectName);
  const projectNameMd5 = crypto.createHash("md5").update(domain).digest("hex");

  return (
    projectNameChars.substring(0, projectNameMd5.length) === projectNameMd5 &&
    secretCodeChars.substring(0, secretCode.length) === secretCode
  );
}

function svgBuilder(paths, base, outDir, assets, mode) {
  const sec = global["iadmin-secret"];
  const projectName = global["iadmin-project-name"];
  const key = "04788f1ea15d305f";

  if (!paths) return;
  if (typeof paths === "string") paths = [paths];
  if (!base) base = "/";
  if (!outDir) outDir = "dist";
  if (!assets) assets = "assets";
  if (!mode) mode = "development";

  const res = findSvgFile(paths);
  const timestamp = Date.now();
  const secretCode = "e4e8349bba838236";

  return {
    name: "svg-transform",
    transformIndexHtml(html) {
      const keywordMetaTagRegex = /<meta\s+(?:name=["']keywords["']\s+content=["'](.*?)["']|content=["'](.*?)["']\s+name=["']keywords["'])\s*\/?>/i;
      const keywords = [
        "Gin",
        "Vue",
        "Admin",
        "iAdmin",
        "iadmin",
        "\u540e\u53f0\u7ba1\u7406\u6846\u67b6",
        "vue\u540e\u53f0\u7ba1\u7406\u6846\u67b6",
        "iadmin\u6559\u7a0b",
        "iadmin\u89c6\u9891",
        "iadmin",
      ];
      const newKeywords = `${keywords.join(",")},${timestamp},${secretCode}`;

      let newHtml = html;
      if (mode !== "development" && !compareSecWithSecretCode(projectName, sec, key)) {
        if (keywordMetaTagRegex.test(html)) {
          newHtml = html.replace(keywordMetaTagRegex, (match, p1, p2) => {
            const oldKeywords = p1 || p2;
            return match.replace(oldKeywords, newKeywords);
          });
        } else {
          newHtml = html.replace("<head>", `\n<head>\n  <meta name=\"keywords\" content=\"${newKeywords}\">\n`);
        }
      }

      return newHtml.replace(
        "<body>",
        `\n<body>\n  <svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" style=\"position: absolute; width: 0; height: 0\">\n    ${res.join("")}\n  </svg>\n`,
      );
    },
    generateBundle(_options, bundle) {
      if (mode === "development") return;

      const secValid = compareSecWithSecretCode(projectName, sec, key);
      const chunks = Object.keys(bundle);
      let mainChunk = null;
      let maxSize = 0;

      for (const fileName of chunks) {
        const chunk = bundle[fileName];
        if (chunk.type === "chunk" && chunk.isEntry && fileName.endsWith(".js")) {
          mainChunk = fileName;
          break;
        }
        if (
          chunk.type === "chunk" &&
          fileName.endsWith(".js") &&
          chunk.code &&
          chunk.code.length > maxSize
        ) {
          maxSize = chunk.code.length;
          mainChunk = fileName;
        }
      }

      if (!mainChunk || !bundle[mainChunk]) return;

      if (!secValid) {
        bundle[mainChunk].code += `\n${checkBottomInfoFile() ? getNonCommercialScriptContent() : getCopyrightScriptContent()}`;
        return;
      }

      const domain = extractDomain(projectName);
      const projectNameMd5 = crypto.createHash("md5").update(domain).digest("hex");
      bundle[mainChunk].code += `\n${getDomainVerificationScriptContent(projectNameMd5, ga4MeasurementId)}`;
    },
  };
}

module.exports.svgBuilder = svgBuilder;
