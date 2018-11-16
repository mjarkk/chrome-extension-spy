import {html, render} from 'lit-html'
import '../style/index.styl'

const log = console.log

log('View the full source code on github: https://github.com/mjarkk/chrome-extension-spy')

const setup = async () => {
  const r = await fetch('http://localhost:8080/lastRequests/')
  const data = await r.json()
  console.log(data)
}

const lastReqests = [
  
]

const addToLastRequests = item => {
  lastReqests.push(item)
  r()
}

const getList = () => html`
  
`

const toRender = () => html`
  <div class="intro-screen flex">
    <h2>Web request</h2>
    <div class="list">
      ${getList()}
    </div>
  </div>
`

const r = () => 
  render(toRender(), document.body)
r()
setup()
