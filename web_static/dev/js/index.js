import {html, render} from 'lit-html'
import '../style/index.styl'

const log = console.log

log('View the full source code on github: https://github.com/mjarkk/chrome-extension-spy')

const setup = () => {
  fetch('/lastRequests')
  .then(r => r.json())
  .then(data => {
    lastReqests = data.reverse()
    dataLoaded = true
    return fetch('/extensionsInfo')
  })
  .then(r => r.json())
  .then(data => {
    extensions = data
    r()
  })
}

const checkIfJson = input => {
  let returnValue = input
  try {
    const testValue = JSON.parse(input)
    returnValue = testValue
  } catch (error) {
    console.log('can\'t convert input to json')
  }
  return returnValue
}

const nicifyOutput = input => {
  const testValue = checkIfJson(input)
  if (typeof testValue == 'object') {
    return JSON.stringify(testValue, null, 2)
  } else {
    return input
  }
}

let extensions = {}
let lastReqests = []
let dataLoaded = false

const extItem = (pkgId, path) => 
  path
    .split('.')
    .reduce((acc, val) => acc && acc[val] ? acc[val] : undefined, extensions[pkgId])

const statusColor = c => 
  c >= 400
    ? 'red'
    : 'green'

const moreIcon = html`
  <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path opacity=".87" fill="none" d="M24 24H0V0h24v24z"/><path d="M15.88 9.29L12 13.17 8.12 9.29c-.39-.39-1.02-.39-1.41 0-.39.39-.39 1.02 0 1.41l4.59 4.59c.39.39 1.02.39 1.41 0l4.59-4.59c.39-.39.39-1.02 0-1.41-.39-.38-1.03-.39-1.42 0z"/></svg>
`

const closeIcon = html`
  <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="none" d="M0 0h24v24H0V0z"/><path d="M18.3 5.71c-.39-.39-1.02-.39-1.41 0L12 10.59 7.11 5.7c-.39-.39-1.02-.39-1.41 0-.39.39-.39 1.02 0 1.41L10.59 12 5.7 16.89c-.39.39-.39 1.02 0 1.41.39.39 1.02.39 1.41 0L12 13.41l4.89 4.89c.39.39 1.02.39 1.41 0 .39-.39.39-1.02 0-1.41L13.41 12l4.89-4.89c.38-.38.38-1.02 0-1.4z"/></svg>
`

const empty_pData = {
  onTab: 0,
  showPopup: false,
  hasLoaded: false,
  firstloadReq: {
    code: 0,
    hash: '',
    pkg: '',
    type: '',
    url: ''
  },
  extension: {
    fullPkgURL: '',
    homepageURL: '',
    name: '',
    pkg: '',
    pkgVersion:'',
    shortName: ''
  },
  req: {}
}

let pData = Object.assign({}, empty_pData)

const loadPopupData = item => {
  pData = Object.assign({}, empty_pData)
  pData.showPopup = true
  pData.firstloadReq = item
  pData.extension = extensions[item.pkg]
  r()
  fetch('/requestInfo/' + item.hash)
  .then(r => r.json())
  .then(data => {
    pData.req = data
    pData.hasLoaded = true
    console.log(pData)
    r()
  })
}

const closePopup = e => {
  if ((e.toElement && e.toElement.classList && e.toElement.classList.value && e.toElement.classList.value.indexOf('popupWrapper') != -1) || (typeof e == 'boolean' && e)) {
    pData.showPopup = false
    r()
  }
}

const popup = () => !pData.showPopup ? html`` : html`
  <div class="popupWrapper flex" @click=${closePopup}>
    <div class="popup">
      <div class="header">
        <div class="row row1">
          <div class="close" @click=${() => closePopup(true)}>${closeIcon}</div>
          <img src="/extLogo/${pData.firstloadReq.pkg}"/>
          <div class="tag">
            <div class="${statusColor(pData.firstloadReq.code)}">${pData.firstloadReq.type} ${pData.firstloadReq.code}</div>
          </div>
        </div>
        <div class="row row2">
          <div 
            class="${pData.onTab == 0 ? 'active' : ''}"
            @click=${() => {
              pData.onTab = 0
              r()
            }}
          >General</div>
          <div 
            class="${pData.onTab == 1 ? 'active' : ''}"
            @click=${() => {
              pData.onTab = 1
              r()
            }}
          >Headers</div>
          <div 
            class="${pData.onTab == 2 ? 'active' : ''}"
            @click=${() => {
              pData.onTab = 2
              r()
            }}
          >Resonse</div>
          <div 
            class="${pData.onTab == 3 ? 'active' : ''}"
            @click=${() => {
              pData.onTab = 3
              r()
            }}
          >PostData</div>
        </div>
      </div>
      ${ pData.onTab == 0 ?
          html`<div class="page page0">
            <div class="info"><span class="item1">Url</span><span class="item2">${pData.firstloadReq.url}</span></div>
            <div class="info"><span class="item1">Type</span><span class="item2">${pData.firstloadReq.type}</span></div>
            <div class="info"><span class="item1">Status</span><span class="item2">${pData.firstloadReq.code}</span></div>
          </div>`
        : pData.onTab == 1 ?
          html`<div class="page page1">
            ${pData.hasLoaded ? html`
              <h3>Request</h3>
              ${pData.req.requestHeaders ? Object.keys(pData.req.requestHeaders).map(el => html`
                <div class="headerItem"><span class="item1">${el}</span><span class="item2">${pData.req.requestHeaders[el]}</span></div> 
              `) : html`<div class="error">The server had an error handling the response headers</div>`}
              <h3>Response</h3>
              ${pData.req.responseHeaders ? Object.keys(pData.req.responseHeaders).map(el => html`
                <div class="headerItem"><span class="item1">${el}</span><span class="item2">${pData.req.responseHeaders[el]}</span></div> 
              `) : html`<div class="error">The server had an error handling the response headers</div>`}
            ` : html`Loading data...`}
          </div>`
        : pData.onTab == 2 ?
          html`<div class="page page2">
            ${pData.hasLoaded ? html`
              <h3>Text data</h3>
              <div class="tip">If the data look wired the response data is probebly not text</div>
              <pre>${ nicifyOutput(pData.req.resData) }</pre>
              <h3>Byte output</h3>
              <pre>${ pData.req.resRawData }</pre>
            ` : html`Loading data...`}
          </div>`
        : pData.onTab == 3 ?
          html`<div class="page page3">
            ${pData.firstloadReq.type != 'POST' 
            ? html `Non post request types don't have post data` : pData.hasLoaded ? html`
              <pre>${nicifyOutput(pData.req.postBody)}</pre>
            ` : html`Loading data...`}
          </div>`
        : html``
      }
    </div>
  </div>
`

const getList = () => html`
  ${lastReqests.length
    ? lastReqests.map(req => html`
        <div class="req">
          <div class="container" @click=${() => loadPopupData(req)}>
            <div class="row row1 flex">
              <div class="logoAndName">
                <img src="/extLogo/${req.pkg}"/>
                <div class="pkgName">${extItem(req.pkg, 'Small.name')}</div>
              </div>
              <div class="moreDetails">${moreIcon}</div>
            </div>
            <div class="row row2 flex">
              <div class="tag"><div class="${statusColor(req.code)}">${req.type} ${req.code}</div></div>
              <div class="url">${req.url}</div>
            </div>
          </div>
        </div>
      `)
    : dataLoaded
      ? html`<div class="no_reqs">
          <h2>Great!</h2>
          <p>Your extensions didn't make any network requests</p>
          <p>Try to visit some sites and see if that changes anything</p>
        </div>`
      : html`<div class="loading">loading data..</div>`
  }
`

const toRender = () => html`
  <div class="intro-screen flex">
    <h2>Web request</h2>
    ${popup()}
    <div class="reqests flex">
      ${getList()}
    </div>
  </div>
`

const r = () => 
  render(toRender(), document.body)
r()
setup()
