import {
  LitElement,
  html,
  css,
} from "https://cdn.jsdelivr.net/gh/lit/dist@3.1.0/core/lit-core.min.js";

class TurnstileComponent extends LitElement {
  static styles = css`
    /* Your CSS styles here */
  `;

  captchaCallback(data) {
    console.log("CALLED!", data);
  }

  async firstUpdated() {
    await this.updateComplete; // Ensures the component's render process is complete
    console.log("in here");

    if (window.turnstile) {
      console.log(this.querySelector("#captcha"));
      turnstile.render(this.querySelector("#captcha"), {
        sitekey: "{{.SiteKey}}",
        callback: function (token) {
          console.log(`Challenge Success ${token}`);
        },
      });
    } else {
      console.error("Turnstile script not loaded");
    }
  }

  createRenderRoot() {
    // Render the template without Shadow DOM. All styles will be global.
    return this;
  }

  render() {
    return html`<section>
      <div id="captcha"></div>
    </section>`;
  }
}

customElements.define("captcha-component", TurnstileComponent);
