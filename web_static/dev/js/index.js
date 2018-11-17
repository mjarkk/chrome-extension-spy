import {html, render} from 'lit-html'
import '../style/index.styl'

const log = console.log

log('View the full source code on github: https://github.com/mjarkk/chrome-extension-spy')

const setup = () => {
  fetch('/lastRequests')
  .then(r => r.json())
  .then(data => {
    lastReqests = data.reverse()
    return fetch('/extensionsInfo')
  })
  .then(r => r.json())
  .then(data => {
    extensions = data
    r()
  })
}

let extensions = {}
let lastReqests = []

const extItem = (pkgId, path) => 
  path
    .split(".")
    .reduce((acc, val) => acc && acc[val] ? acc[val] : undefined, extensions[pkgId])

const statusColor = c => 
  c >= 400
    ? 'red'
    : 'green'

const moreIcon = html`
  <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path opacity=".87" fill="none" d="M24 24H0V0h24v24z"/><path d="M15.88 9.29L12 13.17 8.12 9.29c-.39-.39-1.02-.39-1.41 0-.39.39-.39 1.02 0 1.41l4.59 4.59c.39.39 1.02.39 1.41 0l4.59-4.59c.39-.39.39-1.02 0-1.41-.39-.38-1.03-.39-1.42 0z"/></svg>
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

const popup = () => !pData.showPopup ? html`` : html`
  <div class="popupWrapper">
    <div class="popup">
      <div class="header">
        <div class="row row1">
          <img src="/extLogo/${pData.firstloadReq.pkg}"/>
          <div class="tag">
            <div class="${statusColor(pData.firstloadReq.code)}">${pData.firstloadReq.type} ${pData.firstloadReq.code}</div>
          </div>
        </div>
        <div class="row row2">
          <div @click=${() => {
            pData.onTab = 0
            r()
          }}>General</div>
          <div @click=${() => {
            pData.onTab = 1
            r()
          }}>Headers</div>
          <div @click=${() => {
            pData.onTab = 2
            r()
          }}>Resonse</div>
          <div @click=${() => {
            pData.onTab = 3
            r()
          }}>PostData</div>
        </div>
      </div>
      ${ pData.onTab == 0 ?
          html`<div class="page page0">
            <div class="info"><span class="item1">Url</span><span class="item2">${pData.firstloadReq.url}</span></div>
            <div class="info"><span class="item1">Type</span><span class="item2">${pData.firstloadReq.url}</span></div>
            <div class="info"><span class="item1">Status</span><span class="item2">${pData.firstloadReq.code}</span></div>
          </div>`
        : pData.onTab == 1 ?
          html`<div class="page page1">
            ${pData.hasLoaded ? html`
              <h3>Request</h3>
              ${Object.keys(pData.req.requestHeaders).map(el => html`
                <div class="headerItem"><span class="item1">${el}</span><span class="item2">${pData.req.requestHeaders[el]}</span></div> 
              `)}
              <h3>Response</h3>
              ${Object.keys(pData.req.responseHeaders).map(el => html`
                <div class="headerItem"><span class="item1">${el}</span><span class="item2">${pData.req.responseHeaders[el]}</span></div> 
              `)}
            ` : html`Loading data...`}
          </div>`
        : pData.onTab == 2 ?
          html`<div class="page page2">
            ${pData.hasLoaded ? html`
              <pre>${ pData.req.resData }</pre>
            ` : html`Loading data...`}
          </div>`
        : pData.onTab == 3 ?
          html`<div class="page page3">
            ${pData.firstloadReq.type != "POST" 
            ? html `Non post request types don't have post data` : pData.hasLoaded ? html`
              <pre>${pData.req.postBody}</pre>
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
    : html`<div class="loading">loading data..</div>`}
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
