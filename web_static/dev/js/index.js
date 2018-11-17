import {html, render} from 'lit-html'
import '../style/index.styl'

const log = console.log

log('View the full source code on github: https://github.com/mjarkk/chrome-extension-spy')

const setup = () => {
  fetch('/lastRequests')
  .then(r => r.json())
  .then(data => {
    lastReqests = data
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

const addToLastRequests = item => {
  lastReqests.push(item)
  r()
}

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

const getList = () => html`
  ${lastReqests.length 
    ? lastReqests.reverse().map(req => html`
        <div class="req">
          <div class="container">
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
    <div class="reqests flex">
      ${getList()}
    </div>
  </div>
`

const r = () => 
  render(toRender(), document.body)
r()
setup()
