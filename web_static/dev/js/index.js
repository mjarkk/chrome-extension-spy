import {html, render} from 'lit-html'
import '../style/index.styl'

const log = console.log

log('View the full source code on github: https://github.com/mjarkk/chrome-extension-spy')

const toRender = () => html`
  <div class="intro-screen flex">
    yay
  </div>
`

const r = () => 
  render(toRender(), document.body)
r()
