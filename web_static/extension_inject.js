// TODO make this sand boxed because this is ran global
const CHROME_EXT_SPY_EXT_ID = '--EXT-APP-ID--'
const CHROME_EXT_SPY_BAK_FETCH = window.fetch
const CHROME_EXT_SPY_REPLACE_URI = uri => {
  const dubbleEncode = input => encodeURIComponent(encodeURIComponent(input))
  const toReturn = /\/\/chrome\-extension/.test(uri)
    ? uri
    : /\/\/.+\/\//.test(`http://localhost:8080/proxy/${CHROME_EXT_SPY_EXT_ID}/${uri}`)
      ? `http://localhost:8080/proxy/${CHROME_EXT_SPY_EXT_ID}/${dubbleEncode(uri)}`
      : `http://localhost:8080/proxy/${CHROME_EXT_SPY_EXT_ID}/${dubbleEncode(location.origin + location.pathname.replace(/\/$/, '') + (uri[0] == '/' ? '' : '/') + uri)}`
  return toReturn
}

window.fetch = (uri, options) => {
  return new Promise((resolve, reject) => {
    console.log('req url:', uri)
    CHROME_EXT_SPY_BAK_FETCH(CHROME_EXT_SPY_REPLACE_URI(uri), options)
    .then(res => {
      console.log('req status:', res.status)
      const returnRes = {
        text: async () => {
          const data = await res.text()
          console.log('req text:', data)
          return data
        },
        json: async () => {
          const data = await res.json()
          console.log('req json:', data)
          return data
        },
        arrayBuffer: () => res.arrayBuffer(),
        blob: () => res.blob(),
        clone: () => res.clone(),
        formData: () => res.formData(),
        get body() { return res.body},
        get ok() { return res.ok},
        get redirected() { return res.redirected},
        get status() {return res.status},
        get statusText() {return res.statusText},
        get type() { return res.type},
        get url() { return res.url}
      }
      resolve(returnRes)
    })
    .catch(reject)
  })
}

(function() {
  var origOpen = XMLHttpRequest.prototype.open
  XMLHttpRequest.prototype.open = function(type, uri) {
    arguments[1] = CHROME_EXT_SPY_REPLACE_URI(arguments[1])
    console.log('req url:', uri)
    this.addEventListener('load',() => {
      console.log('req status:', this.readyState)
      console.log('req text:', this.responseText)
    });
    origOpen.apply(this, arguments);
  };
})();
