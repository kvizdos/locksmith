/**
 * @license
 * Copyright 2019 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */
const t=globalThis,e=t.ShadowRoot&&(void 0===t.ShadyCSS||t.ShadyCSS.nativeShadow)&&"adoptedStyleSheets"in Document.prototype&&"replace"in CSSStyleSheet.prototype,i=Symbol(),o=new WeakMap;let r=class{constructor(t,e,o){if(this._$cssResult$=!0,o!==i)throw Error("CSSResult is not constructable. Use `unsafeCSS` or `css` instead.");this.cssText=t,this.t=e}get styleSheet(){let t=this.o;const i=this.t;if(e&&void 0===t){const e=void 0!==i&&1===i.length;e&&(t=o.get(i)),void 0===t&&((this.o=t=new CSSStyleSheet).replaceSync(this.cssText),e&&o.set(i,t))}return t}toString(){return this.cssText}};const s=(t,...e)=>{const o=1===t.length?t[0]:e.reduce(((e,i,o)=>e+(t=>{if(!0===t._$cssResult$)return t.cssText;if("number"==typeof t)return t;throw Error("Value passed to 'css' function must be a 'css' function result: "+t+". Use 'unsafeCSS' to pass non-literal values, but take care to ensure page security.")})(i)+t[o+1]),t[0]);return new r(o,t,i)},a=e?t=>t:t=>t instanceof CSSStyleSheet?(t=>{let e="";for(const i of t.cssRules)e+=i.cssText;return(t=>new r("string"==typeof t?t:t+"",void 0,i))(e)})(t):t
/**
 * @license
 * Copyright 2017 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */,{is:n,defineProperty:l,getOwnPropertyDescriptor:d,getOwnPropertyNames:c,getOwnPropertySymbols:h,getPrototypeOf:p}=Object,u=globalThis,v=u.trustedTypes,g=v?v.emptyScript:"",f=u.reactiveElementPolyfillSupport,m=(t,e)=>t,w={toAttribute(t,e){switch(e){case Boolean:t=t?g:null;break;case Object:case Array:t=null==t?t:JSON.stringify(t)}return t},fromAttribute(t,e){let i=t;switch(e){case Boolean:i=null!==t;break;case Number:i=null===t?null:Number(t);break;case Object:case Array:try{i=JSON.parse(t)}catch(t){i=null}}return i}},y=(t,e)=>!n(t,e),b={attribute:!0,type:String,converter:w,reflect:!1,hasChanged:y};Symbol.metadata??=Symbol("metadata"),u.litPropertyMetadata??=new WeakMap;class $ extends HTMLElement{static addInitializer(t){this._$Ei(),(this.l??=[]).push(t)}static get observedAttributes(){return this.finalize(),this._$Eh&&[...this._$Eh.keys()]}static createProperty(t,e=b){if(e.state&&(e.attribute=!1),this._$Ei(),this.elementProperties.set(t,e),!e.noAccessor){const i=Symbol(),o=this.getPropertyDescriptor(t,i,e);void 0!==o&&l(this.prototype,t,o)}}static getPropertyDescriptor(t,e,i){const{get:o,set:r}=d(this.prototype,t)??{get(){return this[e]},set(t){this[e]=t}};return{get(){return o?.call(this)},set(e){const s=o?.call(this);r.call(this,e),this.requestUpdate(t,s,i)},configurable:!0,enumerable:!0}}static getPropertyOptions(t){return this.elementProperties.get(t)??b}static _$Ei(){if(this.hasOwnProperty(m("elementProperties")))return;const t=p(this);t.finalize(),void 0!==t.l&&(this.l=[...t.l]),this.elementProperties=new Map(t.elementProperties)}static finalize(){if(this.hasOwnProperty(m("finalized")))return;if(this.finalized=!0,this._$Ei(),this.hasOwnProperty(m("properties"))){const t=this.properties,e=[...c(t),...h(t)];for(const i of e)this.createProperty(i,t[i])}const t=this[Symbol.metadata];if(null!==t){const e=litPropertyMetadata.get(t);if(void 0!==e)for(const[t,i]of e)this.elementProperties.set(t,i)}this._$Eh=new Map;for(const[t,e]of this.elementProperties){const i=this._$Eu(t,e);void 0!==i&&this._$Eh.set(i,t)}this.elementStyles=this.finalizeStyles(this.styles)}static finalizeStyles(t){const e=[];if(Array.isArray(t)){const i=new Set(t.flat(1/0).reverse());for(const t of i)e.unshift(a(t))}else void 0!==t&&e.push(a(t));return e}static _$Eu(t,e){const i=e.attribute;return!1===i?void 0:"string"==typeof i?i:"string"==typeof t?t.toLowerCase():void 0}constructor(){super(),this._$Ep=void 0,this.isUpdatePending=!1,this.hasUpdated=!1,this._$Em=null,this._$Ev()}_$Ev(){this._$ES=new Promise((t=>this.enableUpdating=t)),this._$AL=new Map,this._$E_(),this.requestUpdate(),this.constructor.l?.forEach((t=>t(this)))}addController(t){(this._$EO??=new Set).add(t),void 0!==this.renderRoot&&this.isConnected&&t.hostConnected?.()}removeController(t){this._$EO?.delete(t)}_$E_(){const t=new Map,e=this.constructor.elementProperties;for(const i of e.keys())this.hasOwnProperty(i)&&(t.set(i,this[i]),delete this[i]);t.size>0&&(this._$Ep=t)}createRenderRoot(){const i=this.shadowRoot??this.attachShadow(this.constructor.shadowRootOptions);return((i,o)=>{if(e)i.adoptedStyleSheets=o.map((t=>t instanceof CSSStyleSheet?t:t.styleSheet));else for(const e of o){const o=document.createElement("style"),r=t.litNonce;void 0!==r&&o.setAttribute("nonce",r),o.textContent=e.cssText,i.appendChild(o)}})(i,this.constructor.elementStyles),i}connectedCallback(){this.renderRoot??=this.createRenderRoot(),this.enableUpdating(!0),this._$EO?.forEach((t=>t.hostConnected?.()))}enableUpdating(t){}disconnectedCallback(){this._$EO?.forEach((t=>t.hostDisconnected?.()))}attributeChangedCallback(t,e,i){this._$AK(t,i)}_$EC(t,e){const i=this.constructor.elementProperties.get(t),o=this.constructor._$Eu(t,i);if(void 0!==o&&!0===i.reflect){const r=(void 0!==i.converter?.toAttribute?i.converter:w).toAttribute(e,i.type);this._$Em=t,null==r?this.removeAttribute(o):this.setAttribute(o,r),this._$Em=null}}_$AK(t,e){const i=this.constructor,o=i._$Eh.get(t);if(void 0!==o&&this._$Em!==o){const t=i.getPropertyOptions(o),r="function"==typeof t.converter?{fromAttribute:t.converter}:void 0!==t.converter?.fromAttribute?t.converter:w;this._$Em=o,this[o]=r.fromAttribute(e,t.type),this._$Em=null}}requestUpdate(t,e,i){if(void 0!==t){if(i??=this.constructor.getPropertyOptions(t),!(i.hasChanged??y)(this[t],e))return;this.P(t,e,i)}!1===this.isUpdatePending&&(this._$ES=this._$ET())}P(t,e,i){this._$AL.has(t)||this._$AL.set(t,e),!0===i.reflect&&this._$Em!==t&&(this._$Ej??=new Set).add(t)}async _$ET(){this.isUpdatePending=!0;try{await this._$ES}catch(t){Promise.reject(t)}const t=this.scheduleUpdate();return null!=t&&await t,!this.isUpdatePending}scheduleUpdate(){return this.performUpdate()}performUpdate(){if(!this.isUpdatePending)return;if(!this.hasUpdated){if(this.renderRoot??=this.createRenderRoot(),this._$Ep){for(const[t,e]of this._$Ep)this[t]=e;this._$Ep=void 0}const t=this.constructor.elementProperties;if(t.size>0)for(const[e,i]of t)!0!==i.wrapped||this._$AL.has(e)||void 0===this[e]||this.P(e,this[e],i)}let t=!1;const e=this._$AL;try{t=this.shouldUpdate(e),t?(this.willUpdate(e),this._$EO?.forEach((t=>t.hostUpdate?.())),this.update(e)):this._$EU()}catch(e){throw t=!1,this._$EU(),e}t&&this._$AE(e)}willUpdate(t){}_$AE(t){this._$EO?.forEach((t=>t.hostUpdated?.())),this.hasUpdated||(this.hasUpdated=!0,this.firstUpdated(t)),this.updated(t)}_$EU(){this._$AL=new Map,this.isUpdatePending=!1}get updateComplete(){return this.getUpdateComplete()}getUpdateComplete(){return this._$ES}shouldUpdate(t){return!0}update(t){this._$Ej&&=this._$Ej.forEach((t=>this._$EC(t,this[t]))),this._$EU()}updated(t){}firstUpdated(t){}}$.elementStyles=[],$.shadowRootOptions={mode:"open"},$[m("elementProperties")]=new Map,$[m("finalized")]=new Map,f?.({ReactiveElement:$}),(u.reactiveElementVersions??=[]).push("2.0.4");
/**
 * @license
 * Copyright 2017 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */
const x=globalThis,A=x.trustedTypes,_=A?A.createPolicy("lit-html",{createHTML:t=>t}):void 0,k="$lit$",P=`lit$${Math.random().toFixed(9).slice(2)}$`,R="?"+P,S=`<${R}>`,E=document,O=()=>E.createComment(""),C=t=>null===t||"object"!=typeof t&&"function"!=typeof t,M=Array.isArray,z="[ \t\n\f\r]",L=/<(?:(!--|\/[^a-zA-Z])|(\/?[a-zA-Z][^>\s]*)|(\/?$))/g,U=/-->/g,N=/>/g,T=RegExp(`>|${z}(?:([^\\s"'>=/]+)(${z}*=${z}*(?:[^ \t\n\f\r"'\`<>=]|("|')|))|$)`,"g"),j=/'/g,H=/"/g,B=/^(?:script|style|textarea|title)$/i,I=(t=>(e,...i)=>({_$litType$:t,strings:e,values:i}))(1),D=Symbol.for("lit-noChange"),V=Symbol.for("lit-nothing"),F=new WeakMap,q=E.createTreeWalker(E,129);function K(t,e){if(!M(t)||!t.hasOwnProperty("raw"))throw Error("invalid template strings array");return void 0!==_?_.createHTML(e):e}const W=(t,e)=>{const i=t.length-1,o=[];let r,s=2===e?"<svg>":3===e?"<math>":"",a=L;for(let e=0;e<i;e++){const i=t[e];let n,l,d=-1,c=0;for(;c<i.length&&(a.lastIndex=c,l=a.exec(i),null!==l);)c=a.lastIndex,a===L?"!--"===l[1]?a=U:void 0!==l[1]?a=N:void 0!==l[2]?(B.test(l[2])&&(r=RegExp("</"+l[2],"g")),a=T):void 0!==l[3]&&(a=T):a===T?">"===l[0]?(a=r??L,d=-1):void 0===l[1]?d=-2:(d=a.lastIndex-l[2].length,n=l[1],a=void 0===l[3]?T:'"'===l[3]?H:j):a===H||a===j?a=T:a===U||a===N?a=L:(a=T,r=void 0);const h=a===T&&t[e+1].startsWith("/>")?" ":"";s+=a===L?i+S:d>=0?(o.push(n),i.slice(0,d)+k+i.slice(d)+P+h):i+P+(-2===d?e:h)}return[K(t,s+(t[i]||"<?>")+(2===e?"</svg>":3===e?"</math>":"")),o]};class Y{constructor({strings:t,_$litType$:e},i){let o;this.parts=[];let r=0,s=0;const a=t.length-1,n=this.parts,[l,d]=W(t,e);if(this.el=Y.createElement(l,i),q.currentNode=this.el.content,2===e||3===e){const t=this.el.content.firstChild;t.replaceWith(...t.childNodes)}for(;null!==(o=q.nextNode())&&n.length<a;){if(1===o.nodeType){if(o.hasAttributes())for(const t of o.getAttributeNames())if(t.endsWith(k)){const e=d[s++],i=o.getAttribute(t).split(P),a=/([.?@])?(.*)/.exec(e);n.push({type:1,index:r,name:a[2],strings:i,ctor:"."===a[1]?X:"?"===a[1]?tt:"@"===a[1]?et:Q}),o.removeAttribute(t)}else t.startsWith(P)&&(n.push({type:6,index:r}),o.removeAttribute(t));if(B.test(o.tagName)){const t=o.textContent.split(P),e=t.length-1;if(e>0){o.textContent=A?A.emptyScript:"";for(let i=0;i<e;i++)o.append(t[i],O()),q.nextNode(),n.push({type:2,index:++r});o.append(t[e],O())}}}else if(8===o.nodeType)if(o.data===R)n.push({type:2,index:r});else{let t=-1;for(;-1!==(t=o.data.indexOf(P,t+1));)n.push({type:7,index:r}),t+=P.length-1}r++}}static createElement(t,e){const i=E.createElement("template");return i.innerHTML=t,i}}function J(t,e,i=t,o){if(e===D)return e;let r=void 0!==o?i._$Co?.[o]:i._$Cl;const s=C(e)?void 0:e._$litDirective$;return r?.constructor!==s&&(r?._$AO?.(!1),void 0===s?r=void 0:(r=new s(t),r._$AT(t,i,o)),void 0!==o?(i._$Co??=[])[o]=r:i._$Cl=r),void 0!==r&&(e=J(t,r._$AS(t,e.values),r,o)),e}class G{constructor(t,e){this._$AV=[],this._$AN=void 0,this._$AD=t,this._$AM=e}get parentNode(){return this._$AM.parentNode}get _$AU(){return this._$AM._$AU}u(t){const{el:{content:e},parts:i}=this._$AD,o=(t?.creationScope??E).importNode(e,!0);q.currentNode=o;let r=q.nextNode(),s=0,a=0,n=i[0];for(;void 0!==n;){if(s===n.index){let e;2===n.type?e=new Z(r,r.nextSibling,this,t):1===n.type?e=new n.ctor(r,n.name,n.strings,this,t):6===n.type&&(e=new it(r,this,t)),this._$AV.push(e),n=i[++a]}s!==n?.index&&(r=q.nextNode(),s++)}return q.currentNode=E,o}p(t){let e=0;for(const i of this._$AV)void 0!==i&&(void 0!==i.strings?(i._$AI(t,i,e),e+=i.strings.length-2):i._$AI(t[e])),e++}}class Z{get _$AU(){return this._$AM?._$AU??this._$Cv}constructor(t,e,i,o){this.type=2,this._$AH=V,this._$AN=void 0,this._$AA=t,this._$AB=e,this._$AM=i,this.options=o,this._$Cv=o?.isConnected??!0}get parentNode(){let t=this._$AA.parentNode;const e=this._$AM;return void 0!==e&&11===t?.nodeType&&(t=e.parentNode),t}get startNode(){return this._$AA}get endNode(){return this._$AB}_$AI(t,e=this){t=J(this,t,e),C(t)?t===V||null==t||""===t?(this._$AH!==V&&this._$AR(),this._$AH=V):t!==this._$AH&&t!==D&&this._(t):void 0!==t._$litType$?this.$(t):void 0!==t.nodeType?this.T(t):(t=>M(t)||"function"==typeof t?.[Symbol.iterator])(t)?this.k(t):this._(t)}O(t){return this._$AA.parentNode.insertBefore(t,this._$AB)}T(t){this._$AH!==t&&(this._$AR(),this._$AH=this.O(t))}_(t){this._$AH!==V&&C(this._$AH)?this._$AA.nextSibling.data=t:this.T(E.createTextNode(t)),this._$AH=t}$(t){const{values:e,_$litType$:i}=t,o="number"==typeof i?this._$AC(t):(void 0===i.el&&(i.el=Y.createElement(K(i.h,i.h[0]),this.options)),i);if(this._$AH?._$AD===o)this._$AH.p(e);else{const t=new G(o,this),i=t.u(this.options);t.p(e),this.T(i),this._$AH=t}}_$AC(t){let e=F.get(t.strings);return void 0===e&&F.set(t.strings,e=new Y(t)),e}k(t){M(this._$AH)||(this._$AH=[],this._$AR());const e=this._$AH;let i,o=0;for(const r of t)o===e.length?e.push(i=new Z(this.O(O()),this.O(O()),this,this.options)):i=e[o],i._$AI(r),o++;o<e.length&&(this._$AR(i&&i._$AB.nextSibling,o),e.length=o)}_$AR(t=this._$AA.nextSibling,e){for(this._$AP?.(!1,!0,e);t&&t!==this._$AB;){const e=t.nextSibling;t.remove(),t=e}}setConnected(t){void 0===this._$AM&&(this._$Cv=t,this._$AP?.(t))}}class Q{get tagName(){return this.element.tagName}get _$AU(){return this._$AM._$AU}constructor(t,e,i,o,r){this.type=1,this._$AH=V,this._$AN=void 0,this.element=t,this.name=e,this._$AM=o,this.options=r,i.length>2||""!==i[0]||""!==i[1]?(this._$AH=Array(i.length-1).fill(new String),this.strings=i):this._$AH=V}_$AI(t,e=this,i,o){const r=this.strings;let s=!1;if(void 0===r)t=J(this,t,e,0),s=!C(t)||t!==this._$AH&&t!==D,s&&(this._$AH=t);else{const o=t;let a,n;for(t=r[0],a=0;a<r.length-1;a++)n=J(this,o[i+a],e,a),n===D&&(n=this._$AH[a]),s||=!C(n)||n!==this._$AH[a],n===V?t=V:t!==V&&(t+=(n??"")+r[a+1]),this._$AH[a]=n}s&&!o&&this.j(t)}j(t){t===V?this.element.removeAttribute(this.name):this.element.setAttribute(this.name,t??"")}}class X extends Q{constructor(){super(...arguments),this.type=3}j(t){this.element[this.name]=t===V?void 0:t}}class tt extends Q{constructor(){super(...arguments),this.type=4}j(t){this.element.toggleAttribute(this.name,!!t&&t!==V)}}class et extends Q{constructor(t,e,i,o,r){super(t,e,i,o,r),this.type=5}_$AI(t,e=this){if((t=J(this,t,e,0)??V)===D)return;const i=this._$AH,o=t===V&&i!==V||t.capture!==i.capture||t.once!==i.once||t.passive!==i.passive,r=t!==V&&(i===V||o);o&&this.element.removeEventListener(this.name,this,i),r&&this.element.addEventListener(this.name,this,t),this._$AH=t}handleEvent(t){"function"==typeof this._$AH?this._$AH.call(this.options?.host??this.element,t):this._$AH.handleEvent(t)}}class it{constructor(t,e,i){this.element=t,this.type=6,this._$AN=void 0,this._$AM=e,this.options=i}get _$AU(){return this._$AM._$AU}_$AI(t){J(this,t)}}const ot=x.litHtmlPolyfillSupport;ot?.(Y,Z),(x.litHtmlVersions??=[]).push("3.2.1");
/**
 * @license
 * Copyright 2017 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */
let rt=class extends ${constructor(){super(...arguments),this.renderOptions={host:this},this._$Do=void 0}createRenderRoot(){const t=super.createRenderRoot();return this.renderOptions.renderBefore??=t.firstChild,t}update(t){const e=this.render();this.hasUpdated||(this.renderOptions.isConnected=this.isConnected),super.update(t),this._$Do=((t,e,i)=>{const o=i?.renderBefore??e;let r=o._$litPart$;if(void 0===r){const t=i?.renderBefore??null;o._$litPart$=r=new Z(e.insertBefore(O(),t),t,void 0,i??{})}return r._$AI(t),r})(e,this.renderRoot,this.renderOptions)}connectedCallback(){super.connectedCallback(),this._$Do?.setConnected(!0)}disconnectedCallback(){super.disconnectedCallback(),this._$Do?.setConnected(!1)}render(){return D}};rt._$litElement$=!0,rt.finalized=!0,globalThis.litElementHydrateSupport?.({LitElement:rt});const st=globalThis.litElementPolyfillSupport;st?.({LitElement:rt}),(globalThis.litElementVersions??=[]).push("4.1.1");
/**
 * @license
 * Copyright 2017 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */
const at=t=>(e,i)=>{void 0!==i?i.addInitializer((()=>{customElements.define(t,e)})):customElements.define(t,e)}
/**
 * @license
 * Copyright 2017 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */,nt={attribute:!0,type:String,converter:w,reflect:!1,hasChanged:y},lt=(t=nt,e,i)=>{const{kind:o,metadata:r}=i;let s=globalThis.litPropertyMetadata.get(r);if(void 0===s&&globalThis.litPropertyMetadata.set(r,s=new Map),s.set(i.name,t),"accessor"===o){const{name:o}=i;return{set(i){const r=e.get.call(this);e.set.call(this,i),this.requestUpdate(o,r,t)},init(e){return void 0!==e&&this.P(o,void 0,t),e}}}if("setter"===o){const{name:o}=i;return function(i){const r=this[o];e.call(this,i),this.requestUpdate(o,r,t)}}throw Error("Unsupported decorator location: "+o)};function dt(t){return(e,i)=>"object"==typeof i?lt(t,e,i):((t,e,i)=>{const o=e.hasOwnProperty(i);return e.constructor.createProperty(i,o?{...t,wrapped:!0}:t),o?Object.getOwnPropertyDescriptor(e,i):void 0})(t,e,i)
/**
 * @license
 * Copyright 2017 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */}function ct(t){return dt({...t,state:!0,attribute:!1})}const ht=s`
  input {
    padding: 0.85rem;
    border-radius: 0.5rem;
    border: 1px solid var(--input-border, var(--gray-300, #bdbdbd));
    font-size: 1rem;
    width: 100%;
  }
  input:-webkit-autofill,
  textarea:-webkit-autofill,
  select:-webkit-autofill {
    border: 0;
    -webkit-text-fill-color: black;
    -webkit-box-shadow: 0 0 0px 1000px 0 inset;
    transition: background-color 5000s ease-in-out 0s;
  }
  .input-container *:user-invalid {
    outline: 1px solid var(--danger-600, #e01e47);
  }

  textarea {
    padding: 0.85rem;
    border-radius: 0.5rem;
    border: 1px solid var(--input-border, var(--gray-300, #bdbdbd));
    font-size: 1rem;
    width: 100%;
    resize: vertical;
  }

  .input-container {
    width: 100%;
    display: flex;
    flex-direction: column;
  }

  .input-container > label {
    font-weight: 600;
    font-size: 0.85rem;
  }

  .input-container > label:not(:has(+ p)) {
    margin-bottom: 0.5rem;
  }

  .input-container > label + p {
    color: #656565;
    font-size: 0.85rem;
    margin: 0.25rem 0 0.5rem 0;
  }

  .input-container > label:has(button) {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  .input-container > label button {
    border: 0;
    background: 0;
    padding: 0;
    margin: 0;
    display: flex;
    color: var(--accent);
    align-items: center;
    font-weight: 600;
    cursor: pointer;
    gap: 0.5rem;
  }
  .input-container > label:has(~ input:required)::after {
    content: "*";
    margin-left: 0.25rem;
    color: var(--danger-400);
  }

  .input-container input,
  .input-container textarea {
    padding: 0.85rem 0.85rem;
    border-radius: 0.5rem;
    border: 1px solid var(--input-border, #bdbdbd);
    width: 100%;
    font-size: 1rem;
  }

  .input-container input:focus {
    outline: 2px solid var(--accent);
  }

  .input-container p#error {
    margin-top: 0.25rem;
    color: var(--danger-600);
    font-size: 0.85rem;
    display: none;
  }

  .input-container :user-invalid + p#error {
    display: inherit;
  }

  .input-container :user-invalid {
    border-radius: 0.5rem;
    outline: 2px solid var(--danger-600);
  }
`
/**
 * @license
 * Copyright 2017 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */,pt=1,ut=2,vt=t=>(...e)=>({_$litDirective$:t,values:e});class gt{constructor(t){}get _$AU(){return this._$AM._$AU}_$AT(t,e,i){this._$Ct=t,this._$AM=e,this._$Ci=i}_$AS(t,e){return this.update(t,e)}update(t,e){return this.render(...e)}}
/**
 * @license
 * Copyright 2018 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */const ft=vt(class extends gt{constructor(t){if(super(t),t.type!==pt||"class"!==t.name||t.strings?.length>2)throw Error("`classMap()` can only be used in the `class` attribute and must be the only part in the attribute.")}render(t){return" "+Object.keys(t).filter((e=>t[e])).join(" ")+" "}update(t,[e]){if(void 0===this.st){this.st=new Set,void 0!==t.strings&&(this.nt=new Set(t.strings.join(" ").split(/\s/).filter((t=>""!==t))));for(const t in e)e[t]&&!this.nt?.has(t)&&this.st.add(t);return this.render(e)}const i=t.element.classList;for(const t of this.st)t in e||(i.remove(t),this.st.delete(t));for(const t in e){const o=!!e[t];o===this.st.has(t)||this.nt?.has(t)||(o?(i.add(t),this.st.add(t)):(i.remove(t),this.st.delete(t)))}return D}});var mt=function(t,e,i,o){var r,s=arguments.length,a=s<3?e:null===o?o=Object.getOwnPropertyDescriptor(e,i):o;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(t,e,i,o);else for(var n=t.length-1;n>=0;n--)(r=t[n])&&(a=(s<3?r(a):s>3?r(e,i,a):r(e,i))||a);return s>3&&a&&Object.defineProperty(e,i,a),a};let wt=class extends rt{constructor(){super(...arguments),this.disabled=!1,this.expectLoad=!1,this.loading=!1,this.loadingText=""}render(){return I`<button
      ?disabled=${this.disabled||this.loading}
      class=${ft({loading:this.loading})}
      @click=${()=>{this.dispatchEvent(new Event("fl-click"))}}
    >
      ${this.expectLoad?I`
            <svg
              viewBox="0 0 100 100"
              fill="none"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                d="M50 88.5C42.3854 88.5 34.9418 86.242 28.6105 82.0116C22.2793 77.7811 17.3446 71.7683 14.4306 64.7333C11.5167 57.6984 10.7542 49.9573 12.2398 42.489C13.7253 35.0208 17.3921 28.1607 22.7764 22.7764C28.1607 17.3921 35.0208 13.7253 42.489 12.2398C49.9573 10.7542 57.6984 11.5167 64.7333 14.4306C71.7683 17.3446 77.7811 22.2793 82.0116 28.6105C86.242 34.9418 88.5 42.3854 88.5 50"
                stroke="var(--fl-button-loader-loop, var(--primary-300, #8ec4ff))"
                stroke-width="15"
              />
              <path
                id="spinner"
                d="M88.5 50C88.5 55.0559 87.5042 60.0623 85.5694 64.7333C83.6346 69.4043 80.7987 73.6486 77.2236 77.2236C73.6486 80.7987 69.4043 83.6346 64.7333 85.5694C60.0623 87.5042 55.0559 88.5 50 88.5"
                stroke="var(--fl-button-loader-spinner, var(--primary-600, #1a5cf4))"
                stroke-width="15"
              />
            </svg>
          `:I``}
      ${this.loading&&""!==this.loadingText?this.loadingText:I`<slot></slot>`}
    </button>`}};wt.styles=s`
    :host {
      --button-bg: var(--fl-button-bg, var(--primary-500, #327eff));
      --button-text: var(--fl-button-text, #fff);
      --border-radius: 0.5rem;
      --padding: 0.5rem 1rem;
      --weight: 500;
    }

    :host(.pill) {
      --border-radius: 5rem;
      --padding: 0.5rem 1.15rem;
    }

    :host(.big) {
      --padding: 0.85rem;
      --weight: 600;
    }

    button {
      background-color: var(--button-bg);
      color: var(--button-text);
      border: 0;
      padding: var(--padding);
      border-radius: var(--border-radius);
      cursor: pointer;
      font-weight: var(--weight);
      transition: 200ms;
      display: inline-flex;
      align-items: center;
      justify-content: center;
      font-size: 1rem;
    }

    :host(.big) button {
      width: 100%;
    }

    :host(.plain) button {
      background-color: unset;
      color: var(--button-bg);
      padding: 0;
    }

    button.loading {
      gap: 0.5rem;
    }

    button[disabled] {
      --button-bg: var(--fl-button-bg-disabled, var(--primary-100, #d9eaff));
      --button-text: var(
        --fl-button-text-disabled,
        var(--primary-500, #327eff)
      );
      cursor: not-allowed;
    }

    :host(.plain) button[disabled] {
      color: var(--fl-button-plain-disabled, var(--primary-400, #59a2ff));
    }

    button:not([disabled]):hover {
      --button-bg: var(--fl-button-bg-hover, var(--primary-600, #1a5cf4));
    }

    button:not([disabled]):active {
      --button-bg: var(--fl-button-bg-active, var(--primary-800, #173ab6));
      transition: 50ms;
    }

    svg {
      width: 0;
      height: 12px;
    }

    button.loading svg {
      width: 12px;
    }

    svg path {
      transform-origin: center;
      animation: spin 1200ms linear infinite;
    }

    @keyframes spin {
      from {
        transform: rotate(0deg);
      }
      to {
        transform: rotate(360deg);
      }
    }
  `,mt([dt()],wt.prototype,"disabled",void 0),mt([dt()],wt.prototype,"expectLoad",void 0),mt([dt()],wt.prototype,"loading",void 0),mt([dt()],wt.prototype,"loadingText",void 0),wt=mt([at("button-component")],wt);
/**
 * @license
 * Copyright 2020 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */
const yt=(t,e)=>{const i=t._$AN;if(void 0===i)return!1;for(const t of i)t._$AO?.(e,!1),yt(t,e);return!0},bt=t=>{let e,i;do{if(void 0===(e=t._$AM))break;i=e._$AN,i.delete(t),t=e}while(0===i?.size)},$t=t=>{for(let e;e=t._$AM;t=e){let i=e._$AN;if(void 0===i)e._$AN=i=new Set;else if(i.has(t))break;i.add(t),_t(e)}};
/**
 * @license
 * Copyright 2017 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */function xt(t){void 0!==this._$AN?(bt(this),this._$AM=t,$t(this)):this._$AM=t}function At(t,e=!1,i=0){const o=this._$AH,r=this._$AN;if(void 0!==r&&0!==r.size)if(e)if(Array.isArray(o))for(let t=i;t<o.length;t++)yt(o[t],!1),bt(o[t]);else null!=o&&(yt(o,!1),bt(o));else yt(this,t)}const _t=t=>{t.type==ut&&(t._$AP??=At,t._$AQ??=xt)};class kt extends gt{constructor(){super(...arguments),this._$AN=void 0}_$AT(t,e,i){super._$AT(t,e,i),$t(this),this.isConnected=t._$AU}_$AO(t,e=!0){t!==this.isConnected&&(this.isConnected=t,t?this.reconnected?.():this.disconnected?.()),e&&(yt(this,t),bt(this))}setValue(t){if((t=>void 0===t.strings)(this._$Ct))this._$Ct._$AI(t,this);else{const e=[...this._$Ct._$AH];e[this._$Ci]=t,this._$Ct._$AI(e,this,0)}}disconnected(){}reconnected(){}}
/**
 * @license
 * Copyright 2020 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */const Pt=()=>new Rt;class Rt{}const St=new WeakMap,Et=vt(class extends kt{render(t){return V}update(t,[e]){const i=e!==this.Y;return i&&void 0!==this.Y&&this.rt(void 0),(i||this.lt!==this.ct)&&(this.Y=e,this.ht=t.options?.host,this.rt(this.ct=t.element)),V}rt(t){if(this.isConnected||(t=void 0),"function"==typeof this.Y){const e=this.ht??globalThis;let i=St.get(e);void 0===i&&(i=new WeakMap,St.set(e,i)),void 0!==i.get(this.Y)&&this.Y.call(this.ht,void 0),i.set(this.Y,t),void 0!==t&&this.Y.call(this.ht,t)}else this.Y.value=t}get lt(){return"function"==typeof this.Y?St.get(this.ht??globalThis)?.get(this.Y):this.Y?.value}disconnected(){this.lt===this.ct&&this.rt(void 0)}reconnected(){this.rt(this.ct)}});async function Ot(t){const e=(new TextEncoder).encode(t),i=await crypto.subtle.digest("SHA-256",e);return Array.from(new Uint8Array(i)).map((t=>t.toString(16).padStart(2,"0"))).join("")}async function Ct(){const t=window.screen.height*window.devicePixelRatio,e=window.screen.width*window.devicePixelRatio,i=window.screen.colorDepth;console.log("Color Depth",i);const o=await Ot(JSON.stringify({screenHeight:t,screenWidth:e,colorDepth:i})),r=Intl.DateTimeFormat().resolvedOptions().timeZone,s=window.navigator.hardwareConcurrency,a=window.navigator.language,n=await(async()=>{const t=document.createElement("canvas");t.width=500,t.height=500,t.style.display="none",document.body.appendChild(t);const e=t.getContext("2d");e.textBaseline="top",e.font="14px 'Arial'",e.textBaseline="alphabetic",e.fillStyle="#f60",e.fillRect(125,1,62,20),e.fillStyle="#069",e.fillText("Hello, world!",2,15),e.fillStyle="rgba(102, 204, 0, 0.7)",e.fillText("Hello, world!",4,17),e.fillText("🤙",100,20),e.fillText("🎉",110,25),e.fillText("🤣",115,30);const i=t.toDataURL();t.remove();return await Ot(i)})(),l=await(async()=>{const t=(()=>{const t=document.createElement("canvas");let e;try{e=t.getContext("webgl")||t.getContext("experimental-webgl")}catch(t){console.error("Failed to get WebGL context: ",t)}return e})();if(!t)return null;const e=t.getExtension("WEBGL_debug_renderer_info");if(e){const i={renderer:t.getParameter(e.UNMASKED_RENDERER_WEBGL),vendor:t.getParameter(e.UNMASKED_VENDOR_WEBGL)};return await Ot(JSON.stringify(i))}return await Ot(JSON.stringify("blank"))})(),d=await Ot(JSON.stringify({touchSupport:"ontouchstart"in window||navigator.maxTouchPoints>0,maxTouchPoints:navigator.maxTouchPoints})),c=(null===navigator||void 0===navigator?void 0:navigator.platform)||"unknown",h=await(async()=>{const t=new OfflineAudioContext(1,44100,44100),e=t.createOscillator();e.type="sine",e.frequency.setValueAtTime(1e3,t.currentTime),e.connect(t.destination),e.start(0);const i=(await t.startRendering()).getChannelData(0),o=new Uint8Array(i.length);for(let t=0;t<i.length;t++)o[t]=Math.floor(255*(.5*i[t]+.5));const r=await crypto.subtle.digest("SHA-256",o),s=Array.from(new Uint8Array(r)).map((t=>t.toString(16).padStart(2,"0"))).join("");return s})();let p={screen:o,timezone:r,hardwareConcurrency:s,deviceMemory:"0",canvas:n,lang:a,webgl:l,touch:d,battery:!1,platform:c,audio:h,userAgent:"",windowSize:null,dnt:null,devices:null};const u=await Ot(navigator.userAgent),v=await Ot(JSON.stringify({height:window.innerHeight,width:window.innerWidth})),g=navigator.doNotTrack||!1,f=await(async()=>{try{return(await navigator.mediaDevices.enumerateDevices()).map((t=>({kind:t.kind,label:t.label,deviceId:t.deviceId,groupId:t.groupId})))}catch(t){return console.log(t),[]}})(),m=await Ot(f.map((t=>Object.values(t).join(":"))).join(";"));return p={...p,userAgent:u,windowSize:v,dnt:g,devices:m},p}const Mt={activity:I`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 37 37"><path fill-rule="evenodd" d="M8 2h21a6 6 0 0 1 6 6v21a6 6 0 0 1-6 6H8a6 6 0 0 1-6-6V8a6 6 0 0 1 6-6M0 8a8 8 0 0 1 8-8h21a8 8 0 0 1 8 8v21a8 8 0 0 1-8 8H8a8 8 0 0 1-8-8zm8.5 1a1.5 1.5 0 1 0 0 3h21a1.5 1.5 0 0 0 0-3zM7 18.5A1.5 1.5 0 0 1 8.5 17h18a1.5 1.5 0 0 1 0 3h-18A1.5 1.5 0 0 1 7 18.5M8.5 25a1.5 1.5 0 0 0 0 3h20a1.5 1.5 0 0 0 0-3z" class="primary" clip-rule="evenodd"/></svg>`,alert:I`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 100 100"><g clip-path="url(#a)"><path stroke-width="5" d="m52.165 17.75 34.641 60c.962 1.667-.24 3.75-2.165 3.75H15.359c-1.924 0-3.127-2.083-2.165-3.75l34.64-60c.963-1.667 3.369-1.667 4.331 0" class="primary-stroke"/><path d="M44.414 40.384A5 5 0 0 1 49.4 35h1.202a5 5 0 0 1 4.985 5.383l-1.114 14.475a4.486 4.486 0 0 1-8.945 0z" class="primary"/><circle cx="50" cy="68" r="5" class="primary"/></g><defs><clipPath id="a"><path fill="#fff" d="M0 0h100v100H0z"/></clipPath></defs></svg>`,check:I`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 17 15"><path stroke-width="2" d="m1 9.5 3.695 3.695a1 1 0 0 0 1.5-.098L15.5 1" class="primary-stroke"/></svg>`,checkmark:I`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 33 33"><circle cx="16.15" cy="16.15" r="16.15" class="primary"/><path stroke="#fff" stroke-width="3" d="m8.604 18.867 3.328 3.328a1 1 0 0 0 1.452-.04L24.3 9.962"/></svg>`,clock:I`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 100 100"><path stroke-width="6" d="M68.5 14.526A39.8 39.8 0 0 0 50 10c-22.091 0-40 17.909-40 40s17.909 40 40 40 40-17.909 40-40c0-8.127-1.336-14.688-5.5-21" class="secondary-stroke"/><path d="M87.255 18.607a5 5 0 1 0-7.071-7.071L45.536 46.184a5 5 0 1 0 7.07 7.07zM24.16 82.33a5 5 0 0 0-8.66-5l-5 8.66a5 5 0 1 0 8.66 5zm51.34 0a5 5 0 1 1 8.66-5l5 8.66a5 5 0 0 1-8.66 5z" class="primary"/></svg>`,cog:I`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 95 95"><path fill-rule="evenodd" d="M43 0a5 5 0 0 0-5 5v8.286c0 1.856-1.237 3.473-2.951 4.185-1.715.712-3.71.432-5.024-.88l-5.86-5.86a5 5 0 0 0-7.07 0l-6.365 6.363a5 5 0 0 0 0 7.072l5.86 5.86c1.313 1.312 1.593 3.308.88 5.023C16.76 36.763 15.143 38 13.287 38H5a5 5 0 0 0-5 5v9a5 5 0 0 0 5 5h8.286c1.856 0 3.473 1.237 4.185 2.951.712 1.715.432 3.71-.88 5.024l-5.86 5.86a5 5 0 0 0 0 7.07l6.363 6.364a5 5 0 0 0 7.072 0l5.86-5.86c1.312-1.312 3.308-1.592 5.023-.88S38 79.858 38 81.714V90a5 5 0 0 0 5 5h9a5 5 0 0 0 5-5v-8.286c0-1.856 1.237-3.473 2.951-4.185 1.715-.712 3.71-.432 5.024.88l5.86 5.86a5 5 0 0 0 7.07 0l6.365-6.363a5 5 0 0 0 0-7.071l-5.86-5.86c-1.313-1.313-1.593-3.308-.88-5.024.71-1.714 2.327-2.951 4.183-2.951H90a5 5 0 0 0 5-5v-9a5 5 0 0 0-5-5h-8.286c-1.856 0-3.473-1.237-4.185-2.951-.712-1.715-.432-3.71.88-5.024l5.86-5.86a5 5 0 0 0 0-7.07l-6.363-6.365a5 5 0 0 0-7.071 0l-5.86 5.86c-1.313 1.313-3.308 1.593-5.024.88C58.237 16.76 57 15.143 57 13.287V5a5 5 0 0 0-5-5zm4 62c8.284 0 15-6.716 15-15s-6.716-15-15-15-15 6.716-15 15 6.716 15 15 15" class="primary" clip-rule="evenodd"/></svg>`,email:I`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 100 100"><path d="M48.19 50.952 7 30a6 6 0 0 1 6-6h74a6 6 0 0 1 6 6L52.765 50.93a5 5 0 0 1-4.574.022" class="primary"/><path fill-rule="evenodd" d="M88 26H12a4 4 0 0 0-4 4v41a4 4 0 0 0 4 4h76a4 4 0 0 0 4-4V30a4 4 0 0 0-4-4m-76-4a8 8 0 0 0-8 8v41a8 8 0 0 0 8 8h76a8 8 0 0 0 8-8V30a8 8 0 0 0-8-8z" class="secondary" clip-rule="evenodd"/></svg>`,flag:I`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 9 13"><path d="M8 5V1l-1.175.294a10 10 0 0 1-5.588-.215L1 1v4l.237.08a10 10 0 0 0 5.588.214z" class="secondary"/><path d="M1 12.5V5m0 0V1l.237.08a10 10 0 0 0 5.588.214L8 1v4l-1.175.294a10 10 0 0 1-5.588-.215z" class="primary-stroke"/></svg>`,home:I`<svg xmlns="http://www.w3.org/2000/svg" class="icon-home" viewBox="0 0 24 24"><path d="M9 22H5a1 1 0 0 1-1-1V11l8-8 8 8v10a1 1 0 0 1-1 1h-4a1 1 0 0 1-1-1v-4a1 1 0 0 0-1-1h-2a1 1 0 0 0-1 1v4a1 1 0 0 1-1 1m3-9a2 2 0 1 0 0-4 2 2 0 0 0 0 4" class="primary"/><path d="m12.01 4.42-8.3 8.3a1 1 0 1 1-1.42-1.41l9.02-9.02a1 1 0 0 1 1.41 0l8.99 9.02a1 1 0 0 1-1.42 1.41z" class="secondary"/></svg>`,info:I`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"><circle cx="12" cy="12" r="11.5" stroke="#fff"/><path stroke-width="2" d="M13.5 18.5V13a1 1 0 0 0-1-1H10m3.5 6.5h-4m4 0h3" class="primary-stroke"/><circle cx="12.5" cy="7" r="2" class="primary"/></svg>`,note:I`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 100 100"><g stroke-width="6" clip-path="url(#a)"><path d="M58.657 3H18C9.716 3 3 9.716 3 18v64c0 8.284 6.716 15 15 15h64c8.284 0 15-6.716 15-15V34.629" class="primary-stroke"/><path d="M48.93 54.861 79.801 3.473a1 1 0 0 1 1.358-.35L92.707 9.79a1 1 0 0 1 .406 1.29l-.049.091L62.38 62.25 42.86 76.275z" class="secondary-stroke"/></g><defs><clipPath id="a"><path fill="#fff" d="M0 0h100v100H0z"/></clipPath></defs></svg>`,"person-group":I`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M12 13a3 3 0 0 1 3-3h4a3 3 0 0 1 3 3v3a1 1 0 0 1-1 1h-8a1 1 0 0 1-1-1 1 1 0 0 1-1 1H3a1 1 0 0 1-1-1v-3a3 3 0 0 1 3-3h4a3 3 0 0 1 3 3M7 9a3 3 0 1 1 0-6 3 3 0 0 1 0 6m10 0a3 3 0 1 1 0-6 3 3 0 0 1 0 6" class="secondary"/><path d="M12 13a3 3 0 1 1 0-6 3 3 0 0 1 0 6m-3 1h6a3 3 0 0 1 3 3v3a1 1 0 0 1-1 1H7a1 1 0 0 1-1-1v-3a3 3 0 0 1 3-3" class="primary"/></svg>`,"person-outline":I`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 30 34"><path fill-rule="evenodd" d="M20 7a5 5 0 1 1-10 0 5 5 0 0 1 10 0m2 0A7 7 0 1 1 8 7a7 7 0 0 1 14 0M.187 21.912C-.438 17.746 2.788 14 7 14l1.694 1.376a10 10 0 0 0 12.612 0L23 14c4.212 0 7.438 3.746 6.813 7.912L28 34H2zm22.38-4.983 1.09-.886a4.89 4.89 0 0 1 4.178 5.572L26.278 32H3.722L2.165 21.615a4.89 4.89 0 0 1 4.178-5.572l1.09.886a12 12 0 0 0 15.134 0" class="primary" clip-rule="evenodd"/></svg>`,person:I`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 30 34"><g clip-path="url(#a)"><path fill-rule="evenodd" d="M20 7a5 5 0 1 1-10 0 5 5 0 0 1 10 0m2 0A7 7 0 1 1 8 7a7 7 0 0 1 14 0M.187 21.912C-.438 17.746 2.788 14 7 14l1.694 1.376a10 10 0 0 0 12.612 0L23 14c4.212 0 7.438 3.746 6.813 7.912L28 34H2z" class="primary" clip-rule="evenodd"/></g><defs><clipPath id="a"><path fill="#fff" d="M0 0h30v34H0z"/></clipPath></defs></svg>`,"phone-disabled":I`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 100 100"><g clip-path="url(#a)"><rect width="31" height="11.499" x="37.69" y="4.483" stroke="#A4A4A4" stroke-width="4" rx="3" transform="rotate(60 37.69 4.483)"/><path stroke="#A4A4A4" stroke-width="4" d="M56.212 88a13 13 0 0 1-9.46-4.082l-.233-.255-14.483-16.209a13 13 0 0 1-2.514-4.191l-.109-.31L20.205 35.7a13 13 0 0 1 1.186-10.876l.196-.315 3.355-5.218 9.737 16.23c.21.348.345.735.4 1.136l.018.174.88 11.26a27 27 0 0 0 12.767 20.893l.719.426 6.43 3.689c.383.22.713.52.965.88l.103.158L65.434 88z"/><rect width="31" height="11.499" x="70.69" y="60.732" stroke="#A4A4A4" stroke-width="4" rx="3" transform="rotate(60 70.69 60.732)"/><circle cx="36.869" cy="60.869" r="19.869" class="primary"/><path fill="#fff" fill-rule="evenodd" d="M30.32 49.486a1 1 0 0 0-1.413 0l-3.683 3.682a1 1 0 0 0 0 1.415l5.908 5.907a1 1 0 0 1 0 1.414l-6.103 6.103a1 1 0 0 0 0 1.414l3.55 3.55a1 1 0 0 0 1.414 0l6.103-6.103a1 1 0 0 1 1.414 0l5.907 5.908a1 1 0 0 0 1.415 0l3.682-3.682a1 1 0 0 0 0-1.415l-5.908-5.907a1 1 0 0 1 0-1.414l6.103-6.103a1 1 0 0 0 0-1.415l-3.55-3.55a1 1 0 0 0-1.413 0l-6.104 6.104a1 1 0 0 1-1.414 0z" clip-rule="evenodd"/></g><defs><clipPath id="a"><path fill="#fff" d="M0 0h100v100H0z"/></clipPath></defs></svg>`,phone:I`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 100 100"><g clip-path="url(#a)"><rect width="35" height="15.499" x="37.422" y="-.249" class="secondary" rx="5" transform="rotate(60 37.422 -.25)"/><path d="M24.13 16.854a1 1 0 0 1 1.698.026l10.566 17.61c.399.664.637 1.412.698 2.184l.88 11.26a25 25 0 0 0 12.486 19.74l6.431 3.689a5 5 0 0 1 1.779 1.73l9.402 15.386A1 1 0 0 1 67.217 90H56.212a15 15 0 0 1-11.185-5.005L30.544 68.787a15 15 0 0 1-3.026-5.193L18.311 36.34a15 15 0 0 1 1.593-12.913z" class="primary"/><rect width="35" height="15.499" x="70.422" y="56" class="secondary" rx="5" transform="rotate(60 70.422 56)"/></g><defs><clipPath id="a"><path fill="#fff" d="M0 0h100v100H0z"/></clipPath></defs></svg>`,pin:I`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 100 100"><path fill="#fff" d="M0 0h100v100H0z"/><path fill-rule="evenodd" d="M34.825 12.39A5 5 0 0 1 39.787 8H60.94a5 5 0 0 1 4.971 4.465l5.939 55.142a5 5 0 0 1-4.971 5.535h-5.264q.036-.456.036-.923c0-2.683-.914-5.153-2.447-7.116A5 5 0 0 0 62.37 59.9l-2.89-26.696a5 5 0 0 0-4.971-4.462H46.4a5 5 0 0 0-4.963 4.386l-3.302 26.697A5 5 0 0 0 41.045 65a11.52 11.52 0 0 0-2.493 8.142h-5.551a5 5 0 0 1-4.963-5.61z" class="primary" clip-rule="evenodd"/><circle cx="49.868" cy="72" r="7" class="secondary"/><rect width="8" height="18" x="46" y="75" class="secondary" rx="3"/></svg>`,search:I`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"><path fill-rule="evenodd" d="M9.92 7.93a5.93 5.93 0 1 1 11.858 0 5.93 5.93 0 0 1-11.859 0M15.848 0a7.93 7.93 0 0 0-6.27 12.785L.293 22.07l1.414 1.414 9.286-9.286A7.93 7.93 0 1 0 15.848 0" class="primary" clip-rule="evenodd"/></svg>`,"sign-out":I`<svg xmlns="http://www.w3.org/2000/svg" class="icon-door-exit" viewBox="0 0 24 24"><path d="M11 4h3a1 1 0 0 1 1 1v3a1 1 0 0 1-2 0V6h-2v12h2v-2a1 1 0 0 1 2 0v3a1 1 0 0 1-1 1h-3v1a1 1 0 0 1-1.27.96l-6.98-2A1 1 0 0 1 2 19V5a1 1 0 0 1 .75-.97l6.98-2A1 1 0 0 1 11 3z" class="primary"/><path d="m18.59 11-1.3-1.3c-.94-.94.47-2.35 1.42-1.4l3 3a1 1 0 0 1 0 1.4l-3 3c-.95.95-2.36-.46-1.42-1.4l1.3-1.3H14a1 1 0 0 1 0-2z" class="secondary"/></svg>`,sort:I`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 76 76"><rect width="70" height="11" x="3" y="16" class="primary" rx="5"/><rect width="62" height="11" x="11" y="33" class="primary" rx="5"/><rect width="54" height="11" x="19" y="50" class="primary" rx="5"/></svg>`,trash:I`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 57 58"><path stroke-width="2" d="M6.13 18.658h44.74a4 4 0 0 1 3.918 4.804l-6.023 29.356a4 4 0 0 1-4.232 3.184L28.97 54.778a6 6 0 0 0-.94 0l-15.563 1.224a4 4 0 0 1-4.232-3.184L2.212 23.462a4 4 0 0 1 3.918-4.805" class="primary-stroke"/><rect width="4" height="30.964" class="secondary" rx="2" transform="matrix(.99209 -.12553 .2006 .97967 9.295 22.952)"/><rect width="4" height="30.964" class="secondary" rx="2" transform="matrix(.9921 .12548 -.20051 .9797 44.157 22.45)"/><rect width="4" height="28.805" x="26.872" y="22.138" class="secondary" rx="2"/><path fill-rule="evenodd" d="M37.036 0a3.68 3.68 0 0 1 3.678 3.679 3.68 3.68 0 0 0 3.679 3.678h9.664a2.943 2.943 0 0 1 0 5.886H2.943a2.943 2.943 0 0 1 0-5.886h9.664a3.68 3.68 0 0 0 3.679-3.678A3.68 3.68 0 0 1 19.964 0zM22.564 2.207a2.207 2.207 0 1 0 0 4.415h11.872a2.207 2.207 0 0 0 0-4.415z" class="primary" clip-rule="evenodd"/></svg>`,unlink:I`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 100 100"><path stroke-width="8" d="m41.035 46-24.49 24.452a5 5 0 0 0 0 7.077l5.957 5.945a5 5 0 0 0 7.065 0L46 67.067M58.195 54l25.103-24.393a5 5 0 0 0-.01-7.183l-6.276-6.06a5 5 0 0 0-6.957.009L53 32.933" class="primary-stroke"/><rect width="8" height="18" x="65" y="74.997" class="shadow" rx="4" transform="rotate(-45 65 74.997)"/><rect width="8" height="18" x="73.498" y="63.489" class="shadow" rx="4" transform="rotate(-75 73.498 63.489)"/><rect width="8" height="18" x="49.681" y="79.357" class="shadow" rx="4" transform="rotate(-15 49.68 79.357)"/><rect width="8" height="18" x="34.445" y="21.543" class="shadow" rx="4" transform="rotate(135 34.445 21.543)"/><rect width="8" height="18" x="24.947" y="33.05" class="shadow" rx="4" transform="rotate(105 24.947 33.05)"/><rect width="8" height="18" x="49.765" y="18.182" class="shadow" rx="4" transform="rotate(165 49.765 18.182)"/></svg>`,"view-hidden":I`<svg xmlns="http://www.w3.org/2000/svg" class="icon-view-hidden" viewBox="0 0 24 24"><path d="M15.1 19.34a8 8 0 0 1-8.86-1.68L1.3 12.7a1 1 0 0 1 0-1.42L4.18 8.4l2.8 2.8a5 5 0 0 0 5.73 5.73l2.4 2.4zM8.84 4.6a8 8 0 0 1 8.7 1.74l4.96 4.95a1 1 0 0 1 0 1.42l-2.78 2.78-2.87-2.87a5 5 0 0 0-5.58-5.58L8.85 4.6z" class="primary"/><path d="m3.3 4.7 16 16a1 1 0 0 0 1.4-1.4l-16-16a1 1 0 0 0-1.4 1.4" class="secondary"/></svg>`,"view-visible":I`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M17.56 17.66a8 8 0 0 1-11.32 0L1.3 12.7a1 1 0 0 1 0-1.42l4.95-4.95a8 8 0 0 1 11.32 0l4.95 4.95a1 1 0 0 1 0 1.42l-4.95 4.95zM11.9 17a5 5 0 1 0 0-10 5 5 0 0 0 0 10" class="primary"/><circle cx="12" cy="12" r="3" class="secondary"/></svg>`,xmark:I`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 33 33"><circle cx="16.149" cy="16.149" r="16.149" class="primary"/><path stroke="#fff" stroke-width="3" d="m9.81 9.96 6.34 6.34m6.339 6.339-6.34-6.339m0 0 6.34-6.34m-6.34 6.34-6.338 6.339"/></svg>`};var zt,Lt=function(t,e,i,o){var r,s=arguments.length,a=s<3?e:null===o?o=Object.getOwnPropertyDescriptor(e,i):o;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(t,e,i,o);else for(var n=t.length-1;n>=0;n--)(r=t[n])&&(a=(s<3?r(a):s>3?r(e,i,a):r(e,i))||a);return s>3&&a&&Object.defineProperty(e,i,a),a};let Ut=zt=class extends rt{constructor(){super(...arguments),this.name="info",this.size="24px",this.hoverable=!1,this.colorway="primary"}render(){var t;const e=null!==(t=zt.colorways[this.colorway])&&void 0!==t?t:zt.colorways.primary;return I`
      <div
        class=${ft({hoverable:this.hoverable})}
        style="
          --size: ${this.size};
          --primary: ${e.primary};
          --secondary: ${e.secondary};
          --shadow: ${e.shadow};
        "
      >
        ${Mt[this.name]}
      </div>
    `}};function Nt(t){return"function"==typeof t?t():t}Ut.colorways={primary:{primary:"var(--primary-600)",secondary:"var(--primary-500, #327eff)",shadow:"var(--gray-400, #989898)"},danger:{primary:"var(--danger-600, red)",secondary:"var(--danger-500, pink)",shadow:"var(--gray-500, #888)"}},Ut.styles=s`
    :host {
      display: inline-block;
    }

    svg {
      width: var(--size);
      height: var(--size);
      display: block;
    }

    .hoverable svg * {
      transition: 200ms;
    }

    svg .primary {
      fill: var(--primary);
    }
    .hoverable:not(:hover) svg .primary {
      fill: var(--shadow);
    }

    svg .primary-stroke {
      stroke: var(--primary);
    }
    .hoverable:not(:hover) svg .primary-stroke {
      stroke: var(--shadow);
    }

    svg .secondary {
      fill: var(--secondary);
    }
    .hoverable:not(:hover) svg .secondary {
      fill: var(--shadow);
    }
    svg .secondary-stroke {
      stroke: var(--secondary);
    }
    .hoverable:not(:hover) svg .secondary-stroke {
      stroke: var(--shadow);
    }

    svg .shadow {
      fill: var(--shadow);
    }
    svg .shadow-stroke {
      stroke: var(--shadow);
    }
  `,Lt([dt()],Ut.prototype,"name",void 0),Lt([dt()],Ut.prototype,"size",void 0),Lt([dt({type:Boolean})],Ut.prototype,"hoverable",void 0),Lt([dt()],Ut.prototype,"colorway",void 0),Ut=zt=Lt([at("ui-icon")],Ut);class Tt extends Event{static{this.eventName="lit-state-changed"}constructor(t,e,i){super(Tt.eventName,{cancelable:!1}),this.key=t,this.value=e,this.state=i}}const jt=(t,e)=>e!==t&&(e==e||t==t);class Ht extends EventTarget{static{this.finalized=!1}static initPropertyMap(){this.propertyMap||(this.propertyMap=new Map)}get propertyMap(){return this.constructor.propertyMap}get stateValue(){return Object.fromEntries([...this.propertyMap].map((([t])=>[t,this[t]])))}constructor(){super(),this.hookMap=new Map,this.constructor.finalize(),this.propertyMap&&[...this.propertyMap].forEach((([t,e])=>{if(void 0!==e.initialValue){const i=Nt(e.initialValue);this[t]=i,e.value=i}}))}static finalize(){if(this.finalized)return!1;this.finalized=!0;const t=Object.keys(this.properties||{});for(const e of t)this.createProperty(e,this.properties[e]);return!0}static createProperty(t,e){this.finalize();const i="symbol"==typeof t?Symbol():`__${t}`,o=this.getPropertyDescriptor(String(t),i,e);Object.defineProperty(this.prototype,t,o)}static getPropertyDescriptor(t,e,i){const o=i?.hasChanged||jt;return{get(){return this[e]},set(i){const r=this[t];this[e]=i,!0===o(i,r)&&this.dispatchStateEvent(t,i,this)},configurable:!0,enumerable:!0}}reset(){this.hookMap.forEach((t=>t.reset())),[...this.propertyMap].filter((([t,e])=>!(!0===e.skipReset||void 0===e.resetValue))).forEach((([t,e])=>{this[t]=e.resetValue}))}subscribe(t,e,i){e&&!Array.isArray(e)&&(e=[e]);const o=i=>{e&&!e.includes(i.key)||t(i.key,i.value,this)};return this.addEventListener(Tt.eventName,o,i),()=>this.removeEventListener(Tt.eventName,o)}dispatchStateEvent(t,e,i){this.dispatchEvent(new Tt(t,e,i))}}class Bt{constructor(t,e,i){this.host=t,this.state=e,this.callback=i||(()=>this.host.requestUpdate()),this.host.addController(this)}hostConnected(){this.state.addEventListener(Tt.eventName,this.callback),this.callback()}hostDisconnected(){this.state.removeEventListener(Tt.eventName,this.callback)}}function It(t){return(e,i)=>{if(Object.getOwnPropertyDescriptor(e,i))throw new Error("@property must be called before all state decorators");const o=e.constructor;o.initPropertyMap();const r=e.hasOwnProperty(i);return o.propertyMap.set(i,{...t,initialValue:t?.value,resetValue:t?.value}),o.createProperty(i,t),r?Object.getOwnPropertyDescriptor(e,i):void 0}}new URL(location.href);const Dt={prefix:"_ls"};function Vt(t){return t={...Dt,...t},(e,i)=>{const o=Object.getOwnPropertyDescriptor(e,i);if(!o)throw new Error("@local-storage decorator need to be called after @property");const r=`${t?.prefix||""}_${t?.key||String(i)}`,s=e.constructor,a=s.propertyMap.get(i),n=a?.type;if(a){const e=a.initialValue;a.initialValue=()=>function(t,e){if(null!==t&&(e===Boolean||e===Number||e===Array||e===Object))try{t=JSON.parse(t)}catch(e){console.warn("cannot parse value",t)}return t}(localStorage.getItem(r),n)??Nt(e),s.propertyMap.set(i,{...a,...t})}const l=o?.set,d={...o,set:function(t){void 0!==t&&localStorage.setItem(r,n===Object||n===Array?JSON.stringify(t):t),l&&l.call(this,t)}};Object.defineProperty(s.prototype,i,d)}}var Ft=function(t,e,i,o){var r,s=arguments.length,a=s<3?e:null===o?o=Object.getOwnPropertyDescriptor(e,i):o;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(t,e,i,o);else for(var n=t.length-1;n>=0;n--)(r=t[n])&&(a=(s<3?r(a):s>3?r(e,i,a):r(e,i))||a);return s>3&&a&&Object.defineProperty(e,i,a),a};class qt extends Ht{get info(){return this.infoRaw?JSON.parse(this.infoRaw):void 0}set info(t){this.infoRaw=t?JSON.stringify(t):""}get username(){return this.info.username.split("@")[0]}get email(){return this.info.email}get role(){return this.info.role}set extras(t){this.rawExtras=JSON.stringify(t)}get extras(){if(""===this.rawExtras||void 0===this.rawExtras)return{};try{return JSON.parse(this.rawExtras)}catch(t){return console.error(t),{}}}async loadIfNeeded(t=!1){var e;const i=Date.now(),o=""!==this.infoRaw,r=(null===(e=this.permissions)||void 0===e?void 0:e.length)>0,s=0!=+this.loadedAt&&i-+this.loadedAt>9e5;if(!o||!r||s||t)try{const t=await fetch("/api/management/me");if(!t.ok)throw new Error("Failed to fetch /me");const e=await t.json();if(this.info=e.info,this.permissions=e.permissions.join(","),this.loadedAt=`${i}`,qt.GetExtraAboutMe){const t=await qt.GetExtraAboutMe({id:e.info.id,username:e.info.username,role:e.info.role,permissions:e.info.permissions});this.extras=t}}catch(t){console.error("Failed to load /me:",t),window.location.href="/login"}}hasPermission(t){return(this.permissions||"").split(",").includes(t)}hasRole(t){return Array.isArray(t)?t.some((t=>this.role===t)):this.role===t}get isLaunchpad(){const t=document.cookie.split(";");for(let e of t)if(e=e.trim(),e.startsWith("LaunchpadUser="))return e.substring(14)}clear(){localStorage.removeItem("_identity_i"),localStorage.removeItem("_identity_p"),localStorage.removeItem("_identity_la"),localStorage.removeItem("_identity_x")}signOut(){void 0!==qt.SignOutCallback&&qt.SignOutCallback({id:this.info.id,username:this.info.username,role:this.info.role,permissions:this.permissions.split(",")}),this.clear(),window.location.href="/sign-out"}}qt.SignOutCallback=void 0,qt.GetExtraAboutMe=void 0,Ft([Vt({key:"i",prefix:"_identity"}),It()],qt.prototype,"infoRaw",void 0),Ft([Vt({key:"x",prefix:"_identity"}),It()],qt.prototype,"rawExtras",void 0),Ft([Vt({key:"p",prefix:"_identity"}),It()],qt.prototype,"permissions",void 0),Ft([Vt({key:"la",prefix:"_identity"}),It()],qt.prototype,"loadedAt",void 0);const Kt=new qt;var Wt=function(t,e,i,o){var r,s=arguments.length,a=s<3?e:null===o?o=Object.getOwnPropertyDescriptor(e,i):o;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(t,e,i,o);else for(var n=t.length-1;n>=0;n--)(r=t[n])&&(a=(s<3?r(a):s>3?r(e,i,a):r(e,i))||a);return s>3&&a&&Object.defineProperty(e,i,a),a};let Yt=class extends rt{constructor(){super(...arguments),this.jsonSettings="",this.appName="FILL ME",this.settings={OauthProviders:[]},this.originOverride="",this.showingPassword=!1,this.errorMsg=void 0,this.loadedOnce=!1,this.isOnboarding=!1,this.signInRef=Pt(),this.emailRef=Pt(),this.passwordRef=Pt(),this.aboutMeState=new Bt(this,Kt)}connectedCallback(){super.connectedCallback();"true"===new URLSearchParams(window.location.search).get("onboard")&&(this.isOnboarding=!0)}firstUpdated(){var t;null===(t=this.emailRef.value)||void 0===t||t.focus(),""!==this.jsonSettings&&setTimeout((()=>{try{this.settings=JSON.parse(this.jsonSettings)}catch(t){throw console.error("Invalid JSON in jsonSettings:",this.jsonSettings),t}this.loadedOnce=!0}))}canSignIn(){var t,e;return(null===(t=this.emailRef.value)||void 0===t?void 0:t.value.length)>0&&(null===(e=this.passwordRef.value)||void 0===e?void 0:e.value.length)>0}lockInputs(){this.emailRef.value.disabled=!0,this.passwordRef.value.disabled=!0}unlockInputs(){this.emailRef.value.disabled=!1,this.passwordRef.value.disabled=!1}async sendLoginRequest(){var t;const e=await Ct(),i={username:this.emailRef.value.value,password:this.passwordRef.value.value,fingerprint:e},o=await fetch(`${null!==(t=this.originOverride)&&void 0!==t?t:""}/api/login`,{method:"POST",body:JSON.stringify(i)});if(200!==o.status){const t=await o.json();if(t.error)throw new Error(t.error);throw new Error("Something went wrong.")}await Kt.loadIfNeeded(!0)}async attemptSignIn(){var t,e,i;if(this.errorMsg=void 0,!this.canSignIn())return this.errorMsg="Please enter your username and password.",void(0===(null===(t=this.emailRef.value)||void 0===t?void 0:t.value.length)?null===(e=this.emailRef.value)||void 0===e||e.focus():null===(i=this.passwordRef.value)||void 0===i||i.focus());this.signInRef.value.loading=!0,this.lockInputs(),this.requestUpdate();try{await this.sendLoginRequest(),window.location.href=this.isOnboarding&&void 0!==this.settings.PathToOnboard?this.settings.PathToOnboard:"/app"}catch(t){console.error(t),this.errorMsg=void 0,await this.updateComplete,this.errorMsg=t.message}finally{this.signInRef.value.loading=!1,this.unlockInputs()}}keydownEvent(t){var e,i;"Enter"===t.key&&(null===(e=this.emailRef.value)||void 0===e?void 0:e.value.length)>0&&(null===(i=this.passwordRef.value)||void 0===i?void 0:i.value.length)>0&&this.attemptSignIn()}render(){return I` <div id="root" class="${this.loadedOnce?"":"hide"}">
      <div id="header">
        <h1>Sign in to ${this.appName}</h1>
        ${!0!==this.settings.PublicRegistrationsDisabled||this.isOnboarding?I`
              <p id="intro">
                ${this.isOnboarding?I`<strong>Thank you for registering.</strong> Please sign
                      in for the first time to ensure everything works.`:I`Need an account?
                      <a href="/register">Create account</a>`}
              </p>
            `:I``}
        <p
          id="error"
          aria-live="assertive"
          role="status"
          aria-atomic="true"
          aria-relevant="additions"
        >
          ${this.errorMsg?this.errorMsg:""}
        </p>
      </div>

      <div id="inputs">
        <div class="input-container">
          <label for="username">Email Address</label>
          <input
            @keydown=${this.keydownEvent}
            id="username"
            ${Et(this.emailRef)}
            autofill="username"
            autocapitalize="off"
            autocapitalize="off"
            placeholder="Your email"
          />
        </div>

        <div class="input-container">
          <label for="password"
            >Password
            <button
              aria-label="${this.showingPassword?"Hide password":"Show password"}"
              tabindex="-1"
              @click=${()=>{const t=this.passwordRef.value.type;this.passwordRef.value.type="password"===t?"text":"password",this.showingPassword="text"!==t}}
            >
              <ui-icon
                name="${this.showingPassword?"view-hidden":"view-visible"}"
                size="1rem"
              ></ui-icon>
              <p>${this.showingPassword?"Hide":"Show"}</p>
            </button>
          </label>
          <input
            @keydown=${this.keydownEvent}
            id="password"
            ${Et(this.passwordRef)}
            autofill="password"
            type="password"
            placeholder="Your password"
          />
        </div>
      </div>

      <div id="signInArea">
        <button-component
          ${Et(this.signInRef)}
          class="big"
          .expectLoad=${!0}
          .loadingText=${"Signing in.."}
          @fl-click=${this.attemptSignIn}
          >Sign in</button-component
        >

        <a href="/reset-password">Forgot Password</a>
      </div>

      ${this.settings.OauthProviders.length>0?I`
            <div class="hr-split">
              <hr />
              <p>Or...</p>
              <hr />
            </div>
          `:void 0}
      ${this.settings.OauthProviders.map((t=>I`
      <a class="oauth" href="/api/auth/oauth/${t}">
        <img src="/api/auth/oauth/${t}/logo"></img>
          Sign in with
          ${t.charAt(0).toUpperCase()+t.slice(1)}
          <span></span></a>
      `))}
    </div>`}};Yt.styles=[ht,s`
      :host {
        display: block;
      }
      * {
        box-sizing: border-box;
        margin: 0;
        touch-action: manipulation;
      }

      #root {
        display: flex;
        gap: 2rem;
        flex-direction: column;
        opacity: 1;
      }

      #root.hide {
        opacity: 0;
      }

      #inputs .input-container:last-of-type {
        margin-top: 1rem;
      }

      #header h1 {
        font-size: 1.5rem;
        font-weight: 500;
      }

      #header p#intro {
        font-weight: 300;
        margin-top: 0.5rem;
        font-size: 0.85rem;
        color: #464646;
      }

      p#error {
        color: #b8123a;
        margin-top: 0.5rem;
      }

      a {
        color: var(--accent);
      }

      a.oauth {
        border: 1px solid var(--input-border, #bdbdbd);
        background: #fff;
        font-size: 1rem;
        padding: 0.75rem;
        border-radius: 0.5rem;
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 1rem;
        color: #000;
        text-decoration: none;
        color: #000;
      }

      a.oauth img {
        height: 1rem;
      }

      button {
        cursor: pointer;
      }
      hr {
        border: 0;
        border-bottom: 1px solid #dcdcdc;
      }
      .hr-split {
        align-items: center;
        display: flex;
        gap: 1rem;
        width: 100%;
      }
      .hr-split hr {
        width: 100%;
        height: 0px;
      }
      .hr-split p {
        font-weight: 600;
        color: #7c7c7c;
      }
      #signInArea a {
        margin-top: 1rem;
        display: block;
        font-size: 0.85rem;
        width: fit-content;
      }

      ui-icon {
        --primary-600: var(--accent);
        --primary-500: var(--accent);
      }
    `],Wt([dt()],Yt.prototype,"jsonSettings",void 0),Wt([dt()],Yt.prototype,"appName",void 0),Wt([dt()],Yt.prototype,"settings",void 0),Wt([dt()],Yt.prototype,"originOverride",void 0),Wt([ct()],Yt.prototype,"showingPassword",void 0),Wt([ct()],Yt.prototype,"errorMsg",void 0),Wt([ct()],Yt.prototype,"loadedOnce",void 0),Wt([ct()],Yt.prototype,"isOnboarding",void 0),Yt=Wt([at("locksmith-login")],Yt);var Jt=function(t,e,i,o){var r,s=arguments.length,a=s<3?e:null===o?o=Object.getOwnPropertyDescriptor(e,i):o;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(t,e,i,o);else for(var n=t.length-1;n>=0;n--)(r=t[n])&&(a=(s<3?r(a):s>3?r(e,i,a):r(e,i))||a);return s>3&&a&&Object.defineProperty(e,i,a),a};let Gt=class extends rt{constructor(){super(...arguments),this.originOverride="",this.forceEmail="",this.inviteCode="",this.appName="FILL ME",this.minimumPasswordLength=6,this.showingPassword=!1,this.errorMsg=void 0,this.signUpRef=Pt(),this.emailRef=Pt(),this.passwordRef=Pt(),this.passwordConfirmationRef=Pt()}firstUpdated(){var t;null===(t=this.emailRef.value)||void 0===t||t.focus()}canSignIn(){var t,e,i;return(null===(t=this.emailRef.value)||void 0===t?void 0:t.value.length)>0&&(null===(e=this.passwordRef.value)||void 0===e?void 0:e.value.length)>0&&(null===(i=this.passwordConfirmationRef.value)||void 0===i?void 0:i.value.length)>0}doPasswordsMatch(){var t,e;return(null===(t=this.passwordRef.value)||void 0===t?void 0:t.value)===(null===(e=this.passwordConfirmationRef.value)||void 0===e?void 0:e.value)}async sendRegistrationRequest(){var t;const e={username:this.emailRef.value.value,email:this.emailRef.value.value,password:this.passwordRef.value.value,code:this.inviteCode},i=await fetch(`${null!==(t=this.originOverride)&&void 0!==t?t:""}/api/register`,{method:"POST",body:JSON.stringify(e)});if(200!==i.status){if(409===i.status)throw new Error("This email is already being used.");if(400===i.status){const t=await i.json();if("password too short"===t.error)throw new Error("Password too short.");if("illegal username characters"===t.error)throw new Error("Email must be a valid email.")}throw new Error("Something went wrong.")}}passwordLongEnough(){var t;return(null===(t=this.passwordRef.value)||void 0===t?void 0:t.value.length)>=this.minimumPasswordLength}async attemptRegistration(){var t,e,i,o,r,s,a;if(this.errorMsg=void 0,!this.canSignIn())return this.errorMsg="Please enter a username and password.",void(0===(null===(t=this.emailRef.value)||void 0===t?void 0:t.value.length)?null===(e=this.emailRef.value)||void 0===e||e.focus():0===(null===(i=this.passwordRef.value)||void 0===i?void 0:i.value.length)?null===(o=this.passwordRef.value)||void 0===o||o.focus():null===(r=this.passwordConfirmationRef.value)||void 0===r||r.focus());if(!this.doPasswordsMatch())return this.errorMsg="The password must match.",void(null===(s=this.passwordConfirmationRef.value)||void 0===s||s.focus());if(!this.passwordLongEnough())return this.errorMsg=`Password must be at least ${this.minimumPasswordLength} characters long.`,void(null===(a=this.passwordRef.value)||void 0===a||a.focus());this.signUpRef.value.loading=!0,this.requestUpdate();try{await this.sendRegistrationRequest(),window.location.href="/login?onboard=true"}catch(t){console.error(t),this.errorMsg=t.message}finally{this.signUpRef.value.loading=!1}}render(){return I` <div id="root">
      <div id="header">
        <h1>Sign up to ${this.appName}</h1>
        ${0===this.forceEmail.length?I`
              <p id="intro">
                Already have an account? <a href="/login">Sign in instead</a>
              </p>
            `:I``}
        <p id="error">${this.errorMsg}</p>
      </div>

      <div id="inputs">
        <div class="input-container">
          <label for="username">Email Address</label>
          <input
            id="username"
            ${Et(this.emailRef)}
            autofill="username"
            autocapitalize="off"
            autocapitalize="off"
            placeholder="Your email"
            value="${this.forceEmail}"
            ?disabled=${this.forceEmail.length>0}
          />
        </div>

        <div class="input-container">
          <label for="password"
            >Password
            <button
              @click=${()=>{this.showingPassword=!this.showingPassword}}
            >
              <ui-icon
                name="${this.showingPassword?"view-hidden":"view-visible"}"
                size="1rem"
              ></ui-icon>
              <p>${this.showingPassword?"Hide":"Show"}</p>
            </button>
          </label>
          <p>Must be at least ${this.minimumPasswordLength} characters long.</p>
          <input
            id="password"
            ${Et(this.passwordRef)}
            autocomplete="new-password"
            type="${this.showingPassword?"text":"password"}"
            placeholder="Your password"
          />
        </div>

        <div class="input-container">
          <label for="password">Confirm your Password</label>
          <input
            id="password"
            ${Et(this.passwordConfirmationRef)}
            autocomplete="new-password"
            type="${this.showingPassword?"text":"password"}"
            placeholder="Confirm your Password"
          />
        </div>
      </div>

      <button-component
        ${Et(this.signUpRef)}
        class="big"
        .expectLoad=${!0}
        .loadingText=${"Signing Up.."}
        @fl-click=${this.attemptRegistration}
        >Sign Up</button-component
      >
    </div>`}};Gt.styles=[ht,s`
      :host {
        display: block;
      }
      * {
        box-sizing: border-box;
        margin: 0;
        touch-action: manipulation;
      }

      #root {
        display: flex;
        gap: 2rem;
        flex-direction: column;
      }

      #inputs {
        display: flex;
        flex-direction: column;
        gap: 1rem;
      }

      #header h1 {
        font-size: 1.5rem;
        font-weight: 500;
      }

      #header p#intro {
        font-weight: 300;
        margin-top: 0.5rem;
        font-size: 0.85rem;
        color: #464646;
        line-height: 1.15rem;
      }

      p#error {
        color: #b8123a;
        margin-top: 0.5rem;
      }

      a {
        color: var(--accent);
      }

      button {
        cursor: pointer;
      }

      ui-icon {
        --primary-600: var(--accent);
        --primary-500: var(--accent);
      }
    `],Jt([dt()],Gt.prototype,"originOverride",void 0),Jt([dt()],Gt.prototype,"forceEmail",void 0),Jt([dt()],Gt.prototype,"inviteCode",void 0),Jt([dt()],Gt.prototype,"appName",void 0),Jt([dt()],Gt.prototype,"minimumPasswordLength",void 0),Jt([ct()],Gt.prototype,"showingPassword",void 0),Jt([ct()],Gt.prototype,"errorMsg",void 0),Gt=Jt([at("locksmith-registration")],Gt);var Zt=function(t,e,i,o){var r,s=arguments.length,a=s<3?e:null===o?o=Object.getOwnPropertyDescriptor(e,i):o;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(t,e,i,o);else for(var n=t.length-1;n>=0;n--)(r=t[n])&&(a=(s<3?r(a):s>3?r(e,i,a):r(e,i))||a);return s>3&&a&&Object.defineProperty(e,i,a),a};let Qt=class extends rt{constructor(){super(...arguments),this.hasResetCode="",this.appName="FILL ME",this.originOverride="",this.minimumPasswordLength=6,this.showingPassword=!1,this.errorMsg=void 0,this.loadedOnce=!1,this.stage=0,this.passwordRef=Pt(),this.passwordConfirmationRef=Pt(),this.resetButtonRef=Pt(),this.emailRef=Pt(),this.resetFullyButtonRef=Pt()}firstUpdated(){var t,e;null===(t=this.emailRef.value)||void 0===t||t.focus(),null===(e=this.passwordRef.value)||void 0===e||e.focus(),setTimeout((()=>{this.loadedOnce=!0}))}async attemptReset(){var t,e;if(0===(null===(t=this.emailRef.value)||void 0===t?void 0:t.value.length))return void(this.errorMsg="Please enter an email address.");this.resetButtonRef.value.loading=!0;if(200!==(await fetch(`${null!==(e=this.originOverride)&&void 0!==e?e:""}/api/reset-password?username=${this.emailRef.value.value}`,{method:"POST"})).status)return console.error("Something bad happened while resetting a password"),void(this.errorMsg="Something went wrong. Please try again later.");this.stage=1}doPasswordsMatch(){var t,e;return(null===(t=this.passwordRef.value)||void 0===t?void 0:t.value)===(null===(e=this.passwordConfirmationRef.value)||void 0===e?void 0:e.value)}passwordLongEnough(){var t;return(null===(t=this.passwordRef.value)||void 0===t?void 0:t.value.length)>=this.minimumPasswordLength}async fullyResetPassword(){var t,e;if(!this.doPasswordsMatch())return void(this.errorMsg="The password must match.");if(!this.passwordLongEnough())return this.errorMsg=`Password must be at least ${this.minimumPasswordLength} characters long.`,void(null===(t=this.passwordRef.value)||void 0===t||t.focus());this.resetFullyButtonRef.value.loading=!0;if(200!==(await fetch(`${null!==(e=this.originOverride)&&void 0!==e?e:""}/api/reset-password`,{method:"PATCH",body:JSON.stringify({password:this.passwordRef.value.value})})).status)return console.error("Something bad happened while resetting a password"),void(this.errorMsg="Something went wrong. Please try again later.");this.stage=2}render(){return 2==this.stage?I`<div id="root" class="${this.loadedOnce?"":"hide"}">
        <div id="header">
          <h1>Your password has been reset.</h1>
          <p id="intro">
            <strong>You may now login using the new password.</strong> If you
            need any help, feel free to contact us.
          </p>
        </div>

        <a href="/login">Go to Login</a>
      </div>`:""!==this.hasResetCode?I`<div id="root" class="${this.loadedOnce?"":"hide"}">
        <div id="header">
          <h1>Welcome back! Please enter a new password</h1>
          <p
            id="error"
            aria-live="assertive"
            role="status"
            aria-atomic="true"
            aria-relevant="additions"
          >
            ${this.errorMsg?this.errorMsg:""}
          </p>
        </div>
        <div id="inputs">
          <div class="input-container">
            <label for="password"
              >Password
              <button
                @click=${()=>{this.showingPassword=!this.showingPassword}}
              >
                <ui-icon
                  name="${this.showingPassword?"view-hidden":"view-visible"}"
                  size="1rem"
                ></ui-icon>
                <p>${this.showingPassword?"Hide":"Show"}</p>
              </button>
            </label>
            <p>
              Must be at least ${this.minimumPasswordLength} characters long.
            </p>
            <input
              id="password"
              ${Et(this.passwordRef)}
              autocomplete="new-password"
              type="${this.showingPassword?"text":"password"}"
              placeholder="Your password"
            />
          </div>

          <div class="input-container">
            <label for="password">Confirm your Password</label>
            <input
              id="password"
              ${Et(this.passwordConfirmationRef)}
              autocomplete="new-password"
              type="${this.showingPassword?"text":"password"}"
              placeholder="Confirm your Password"
            />
          </div>
        </div>
        <button-component
          ${Et(this.resetFullyButtonRef)}
          class="big"
          .expectLoad=${!0}
          .loadingText=${"Resetting.."}
          @fl-click=${this.fullyResetPassword}
          >Reset Password</button-component
        >
      </div>`:1==this.stage?I` <div id="root" class="${this.loadedOnce?"":"hide"}">
        <div id="header">
          <h1>Please check your email.</h1>
          <p id="intro">
            <strong>You may now safely close this window.</strong> If we have an
            account associated with the email address provided, you'll receive
            an email shortly with instructions to reset your password. Be sure
            to check your spam folder. If you need any help, feel free to
            contact us.
          </p>
        </div>

        <a href="/login">Back to Login</a>
      </div>`:I` <div id="root" class="${this.loadedOnce?"":"hide"}">
      <div id="header">
        <h1>Forgot Password</h1>
        <p id="intro">
          Please enter your email address. We will email you a link to reset
          your password.
        </p>
        <p
          id="error"
          aria-live="assertive"
          role="status"
          aria-atomic="true"
          aria-relevant="additions"
        >
          ${this.errorMsg?this.errorMsg:""}
        </p>
      </div>

      <div id="inputs">
        <div class="input-container">
          <label for="username">Email Address</label>
          <input
            id="username"
            ${Et(this.emailRef)}
            autofill="username"
            autocapitalize="off"
            autocapitalize="off"
            placeholder="Your email"
          />
        </div>
      </div>

      <button-component
        ${Et(this.resetButtonRef)}
        class="big"
        .expectLoad=${!0}
        .loadingText=${"Sending.."}
        @fl-click=${this.attemptReset}
        >Send Reset Link</button-component
      >
    </div>`}};Qt.styles=[ht,s`
      :host {
        display: block;
      }
      * {
        box-sizing: border-box;
        margin: 0;
        touch-action: manipulation;
      }

      #root {
        display: flex;
        gap: 2rem;
        flex-direction: column;
        opacity: 1;
      }

      #root.hide {
        opacity: 0;
      }

      #header h1 {
        font-size: 1.5rem;
        font-weight: 500;
      }

      #header p#intro {
        font-weight: 300;
        margin-top: 0.5rem;
        font-size: 0.85rem;
        color: #464646;
        line-height: 1.15rem;
      }

      p#error {
        color: #b8123a;
        margin-top: 0.5rem;
      }

      #inputs {
        display: flex;
        flex-direction: column;
        gap: 1rem;
      }

      a {
        color: var(--accent);
      }

      button {
        cursor: pointer;
      }

      ui-icon {
        --primary-600: var(--accent);
        --primary-500: var(--accent);
      }
    `],Zt([dt()],Qt.prototype,"hasResetCode",void 0),Zt([dt()],Qt.prototype,"appName",void 0),Zt([dt()],Qt.prototype,"originOverride",void 0),Zt([dt()],Qt.prototype,"minimumPasswordLength",void 0),Zt([ct()],Qt.prototype,"showingPassword",void 0),Zt([ct()],Qt.prototype,"errorMsg",void 0),Zt([ct()],Qt.prototype,"loadedOnce",void 0),Zt([ct()],Qt.prototype,"stage",void 0),Qt=Zt([at("locksmith-reset-password")],Qt);var Xt=function(t,e,i,o){var r,s=arguments.length,a=s<3?e:null===o?o=Object.getOwnPropertyDescriptor(e,i):o;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(t,e,i,o);else for(var n=t.length-1;n>=0;n--)(r=t[n])&&(a=(s<3?r(a):s>3?r(e,i,a):r(e,i))||a);return s>3&&a&&Object.defineProperty(e,i,a),a};let te=class extends rt{constructor(){super(...arguments),this.pages=[{URLKey:"account",Name:"Account",TemplateLiteral:I`<p>Hello, Account page</p>`},{URLKey:"security",Name:"Security",TemplateLiteral:I`<p>Hello, Security page</p>`},{URLKey:"demo",Name:"Demo",TemplateLiteral:I`<p>Hello, Demo page</p>`}],this.defaultPageKey="demo",this.selected="",this.activePageComponent=void 0}async setSelected(t,e=!0){var i;if(this.selected=t.URLKey,e&&history.pushState({key:t.URLKey},"",`#${t.URLKey}`),void 0!==t.PageComponent){const e=new t.PageComponent;if(void 0!==t.LoadProps){const i=t.LoadProps();e.setProps(i)}return e.OnPageLoad&&await e.OnPageLoad(),void(this.activePageComponent=e)}this.activePageComponent=null!==(i=t.TemplateLiteral)&&void 0!==i?i:I`<p>Page is missing it's definition.</p.>`}firstUpdated(){var t;const e=(null===(t=location.hash)||void 0===t?void 0:t.replace("#",""))||this.defaultPageKey;let i=this.pages.find((t=>t.URLKey===e));i||(i=this.pages.find((t=>t.URLKey===this.defaultPageKey))),i?this.setSelected(i,!1).then((()=>{this.updateComplete.then((()=>{const t=this.renderRoot.querySelector("button.selected");null==t||t.scrollIntoView({behavior:"smooth",inline:"center",block:"nearest"})}))})):console.warn("Page key not found: ",this.defaultPageKey)}connectedCallback(){super.connectedCallback(),window.addEventListener("popstate",this.listenForNavChanges.bind(this))}disconnectedCallback(){super.connectedCallback(),window.removeEventListener("popstate",this.listenForNavChanges.bind(this))}listenForNavChanges(){const t=location.hash.replace("#","")||this.defaultPageKey,e=this.pages.find((e=>e.URLKey===t));e&&(this.setSelected(e,!1),this.updateComplete.then((()=>{const t=this.renderRoot.querySelector("button.selected");null==t||t.scrollIntoView({behavior:"smooth",inline:"center",block:"nearest"})})))}render(){return I`<nav>
        ${this.pages.map((t=>I`<button
              class=${ft({selected:this.selected===t.URLKey})}
              @click=${()=>this.setSelected(t)}
            >
              ${t.Name}
            </button>`))}
      </nav>
      <section>
        ${void 0!==this.activePageComponent?this.activePageComponent:""}
      </section>`}};te.styles=s`
    :host {
      --hover: var(--quick-nav-hover, var(--gray-50, #f8f8f8));
      --active: var(--quick-nav-active, var(--primary-600, #1a5cf4));
      --focus: var(--quick-nav-focus, var(--primary-500, #327eff));
      --border-radius: 0.5rem;
    }

    nav {
      display: flex;
      gap: 0.25rem;
      position: relative;
      overflow-x: auto;
      -webkit-overflow-scrolling: touch;
      scroll-behavior: smooth;
      scrollbar-width: none; /* Firefox */
      padding: 2px; /* add horizontal padding */
      scroll-padding: 0 0.5rem; /* ensure focused items aren't flush with edges */
    }

    nav::-webkit-scrollbar {
      display: none; /* Chrome/Safari */
    }

    nav button {
      color: #1a1a1a;
      background: unset;
      border: 0;
      padding: 0.5rem 1rem;
      font-size: 1rem;
      transition: background-color 200ms;
      cursor: pointer;
      border-radius: var(--border-radius);
      position: relative;
      min-width: max-content;
      white-space: nowrap;
      flex-shrink: 0;
    }

    nav button:not(.selected):hover {
      background-color: var(--hover);
    }

    nav button:focus-visible {
      background-color: var(--hover);
      outline: 2px solid var(--focus);
    }

    nav button::after {
      content: "";
      position: absolute;
      left: 0;
      bottom: 0;
      height: 2px;
      background: var(--active);
      width: 0%;
      transition: width 200ms ease;
    }

    nav button.selected::after {
      width: 100%;
    }

    section {
      margin-top: 1rem;
    }
  `,Xt([dt()],te.prototype,"pages",void 0),Xt([dt()],te.prototype,"defaultPageKey",void 0),Xt([ct()],te.prototype,"selected",void 0),Xt([ct()],te.prototype,"activePageComponent",void 0),te=Xt([at("quick-nav")],te);var ee=function(t,e,i,o){var r,s=arguments.length,a=s<3?e:null===o?o=Object.getOwnPropertyDescriptor(e,i):o;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(t,e,i,o);else for(var n=t.length-1;n>=0;n--)(r=t[n])&&(a=(s<3?r(a):s>3?r(e,i,a):r(e,i))||a);return s>3&&a&&Object.defineProperty(e,i,a),a};let ie=class extends rt{constructor(){super(...arguments),this.aboutMeState=new Bt(this,Kt)}render(){return I`<div class="input-container">
      <label>
        Email Address
        ${Kt.hasPermission("user.update.email")?I` <button>Change Email</button> `:I``}
      </label>
      ${Kt.hasPermission("user.update.email")?I``:I` <p>Please contact us to change your account email.</p> `}
      <input value="${Kt.email}" disabled />
    </div>`}};ie.styles=[ht,s`
      * {
        box-sizing: border-box;
      }
    `],ie=ee([at("locksmith-update-email")],ie);var oe=function(t,e,i,o){var r,s=arguments.length,a=s<3?e:null===o?o=Object.getOwnPropertyDescriptor(e,i):o;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(t,e,i,o);else for(var n=t.length-1;n>=0;n--)(r=t[n])&&(a=(s<3?r(a):s>3?r(e,i,a):r(e,i))||a);return s>3&&a&&Object.defineProperty(e,i,a),a};let re=class extends rt{constructor(){super(...arguments),this.aboutMeState=new Bt(this,Kt)}connectedCallback(){super.connectedCallback(),Kt.loadIfNeeded(!0)}render(){return I`<div>
      <a href="/app" id="back"
        ><ui-icon name="home" size="1.15rem"></ui-icon> Back to App</a
      >

      <div id="header">
        <h1>Hi, ${Kt.username}.</h1>
        <p id="intro">Manage your account settings here.</p>
        <p
          id="error"
          aria-live="assertive"
          role="status"
          aria-atomic="true"
          aria-relevant="additions"
        ></p>
      </div>

      <quick-nav
        .defaultPageKey=${"account"}
        .pages=${[{URLKey:"account",Name:"Account",PageComponent:ie},{URLKey:"security",Name:"Security",TemplateLiteral:I`<p>TODO</p>`}]}
      ></quick-nav>
    </div>`}};re.styles=[ht,s`
      #header h1 {
        font-size: 1.5rem;
        font-weight: 500;
      }

      #header p#intro {
        font-weight: 300;
        margin-top: 0.5rem;
        font-size: 0.85rem;
        color: #464646;
      }

      a#back {
        gap: 0.5rem;
        color: var(--primary-600, #1a5cf4);
        display: flex;
        align-items: center;
        text-decoration: none;
        font-weight: 300;
      }
    `],re=oe([at("locksmith-profile")],re);var se=function(t,e,i,o){var r,s=arguments.length,a=s<3?e:null===o?o=Object.getOwnPropertyDescriptor(e,i):o;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(t,e,i,o);else for(var n=t.length-1;n>=0;n--)(r=t[n])&&(a=(s<3?r(a):s>3?r(e,i,a):r(e,i))||a);return s>3&&a&&Object.defineProperty(e,i,a),a};let ae=class extends rt{render(){return I`<div id="root">
      <header>
        ${void 0!==this.logoURL&&""!==this.logoURL?I` <img src="${this.logoURL}" /> `:I``}
      </header>
      <main>
        <div id="slotWrapper">
          <slot name="main"></slot>
        </div>
      </main>
      <footer></footer>
    </div>`}};async function ne(t,e){return fetch(t,e)}ae.styles=[s`
      * {
        box-sizing: border-box;
        touch-action: manipulation;
        margin: 0;
      }
      #root {
        display: flex;
        flex-direction: column;
        height: 100svh;
        --horizontal-padding: 1.5rem;
      }

      header,
      footer {
        padding: 1rem 1.5rem;
      }

      footer {
        padding-bottom: 1.5rem;
      }

      header {
        --accent-height: 0.5rem;
        padding-top: calc(1.5rem + var(--accent-height));
      }
      header img {
        height: 2.5rem;
      }

      header::before {
        z-index: -1;
        content: " ";
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: var(--accent-height);
        background-color: var(--accent, #8c7ffa);
      }

      main {
        height: 100%;
        display: flex;
        justify-content: center;
      }

      main #slotWrapper {
        border-radius: 0.35rem;
        padding: 0rem var(--horizontal-padding);
        width: 100%;
      }

      @media (min-width: 650px) {
        #root {
          justify-content: space-between;
          --horizontal-padding: 3rem;
        }
        header {
          padding-bottom: 0;
        }
        main {
          height: 100%;
          align-items: center;
        }
        main #slotWrapper {
          background-color: #fff;
          border: 1px solid #dcdcdc;
          max-width: 28rem;
          padding: 3.5rem var(--horizontal-padding);
        }
      }
    `],se([dt()],ae.prototype,"logoURL",void 0),ae=se([at("locksmith-layout")],ae);var le=function(t,e,i,o){var r,s=arguments.length,a=s<3?e:null===o?o=Object.getOwnPropertyDescriptor(e,i):o;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(t,e,i,o);else for(var n=t.length-1;n>=0;n--)(r=t[n])&&(a=(s<3?r(a):s>3?r(e,i,a):r(e,i))||a);return s>3&&a&&Object.defineProperty(e,i,a),a};let de=class extends rt{constructor(){super(...arguments),this.location={vertical:"top",horizontal:"right"},this.open=!1,this.launchpadForceClosed=!1,this.aboutMeState=new Bt(this,Kt),this.manageAccountButton=Pt(),this.handleEscape=t=>{"Escape"===t.key&&(this.open=!1)},this.listenForOOBClicks=t=>{if(!this.shadowRoot)return;t.composedPath().includes(this.shadowRoot.host)||(this.open=!1,window.removeEventListener("click",this.listenForOOBClicks))}}updated(){this.setAttribute("location-vertical",this.location.vertical),this.setAttribute("location-horizontal",this.location.horizontal)}connectedCallback(){super.connectedCallback(),Kt.loadIfNeeded()}openClicked(){this.open=!this.open,setTimeout((()=>{this.open?(window.addEventListener("click",this.listenForOOBClicks),window.addEventListener("keydown",this.handleEscape),this.updateComplete.then((()=>{setTimeout((()=>{var t;null===(t=this.manageAccountButton.value)||void 0===t||t.focus()}),100)}))):(window.removeEventListener("click",this.listenForOOBClicks),window.removeEventListener("keydown",this.handleEscape))}))}render(){return I` <div id="container">
        <button
          id="box"
          @click=${this.openClicked}
          aria-label="User profile"
          aria-expanded="${this.open?"true":"false"}"
          aria-controls="dropdown"
          aria-haspopup="menu"
        >
          <ui-icon name="person" size="1.25rem"></ui-icon>
        </button>

        <div id="dropdown" class="${this.open?"":"closed"}">
          <div id="header">
            <div>
              <p id="title">${Kt.username}</p>
              <p id="desc">${Kt.email}</p>
            </div>
          </div>

          <div id="actions">
            <button
              ${Et(this.manageAccountButton)}
              @click=${()=>{window.location.href="/profile"}}
            >
              <ui-icon name="cog" size="1rem"></ui-icon>
              Manage Account
              <span></span>
            </button>
            ${Kt.hasPermission("view.ls-admin")?I`
                  <button
                    @click=${()=>{window.location.href="/locksmith"}}
                  >
                    <ui-icon name="person-group" size="1rem"></ui-icon>
                    Administration
                    <span></span>
                  </button>
                `:void 0}
            <button
              id="logout"
              @click=${()=>{Kt.signOut()}}
            >
              <ui-icon name="sign-out" size="1rem" colorway="danger"></ui-icon>
              Sign out
            </button>
          </div>
        </div>
      </div>

      ${void 0===Kt.isLaunchpad||this.launchpadForceClosed?void 0:I`<div id="launchpad">
            <div>
              <p id="launchpad-status">Launchpad</p>
              <p id="launchpad-user">Viewing app as ${Kt.isLaunchpad}</p>
            </div>

            <div>
              <button
                @click=${()=>{this.launchpadForceClosed=!0}}
              >
                &times;
              </button>
            </div>
          </div>`}`}};de.styles=[s`
      #box {
        background-color: #d9eaff;
        position: relative;
        width: 1.25rem;
        aspect-ratio: 1/1;
        overflow: hidden;
        border-radius: 100%;
        border: 2px solid var(--primary-700, #1448e1);
        transition: 200ms;
        outline: 0px solid #8ec4ff;
        cursor: pointer;
        box-sizing: content-box;
        padding: 0;
      }

      #box:hover,
      #box:focus-visible {
        --primary-600: #1a5cf4;
        outline: 3px solid #d9eaff;
        transition: 200ms;
      }

      #box ui-icon {
        position: absolute;
        bottom: -4px;
        left: 0;
        --primary-600: var(--primary-900, #1a5cf4);
      }

      #container {
        display: flex;
        height: 2rem;
        align-items: center;
        gap: 0.85rem;
        position: relative;
      }

      #dropdown {
        margin-top: 0.65rem;
        position: absolute;
        top: 100%;
        border-top: 2px solid var(--primary-700, #1448e1);
        border-right: 1px solid var(--gray-100, #eaeaea);
        border-left: 1px solid var(--gray-100, #eaeaea);
        border-bottom: 1px solid var(--gray-100, #eaeaea);
        padding: 1.15rem;
        border-radius: 0 0 0 0.85rem;
        background-color: #fff;
        --width: 14rem;
        width: var(--width);
        left: calc(-1 * var(--width) - 1rem);
        box-shadow:
          0 4px 12px rgba(0, 0, 0, 0.1),
          0 2px 4px rgba(0, 0, 0, 0.06);
        transition: 200ms;
        opacity: 1;
      }

      #dropdown.closed {
        transform: scale3d(0.9, 0.9, 0.9);
        transition: 200ms;
        opacity: 0;
        visibility: hidden;
        pointer-events: none;
      }

      #dropdown * {
        box-sizing: border-box;
        margin: 0;
      }

      #dropdown #header {
        display: flex;
        gap: 1rem;
        align-items: center;
        --primary-600: var(--primary-800, #173ab6);
        --primary-500: var(--primary-300, #8ec4ff);
      }

      #dropdown #header #title {
        font-size: 1.1rem;
        font-weight: 500;
      }

      #dropdown #header #desc {
        font-size: 0.85rem;
        font-weight: 400;
        margin-top: 0.25rem;
        color: var(--gray-600, #656565);
        text-overflow: ellipsis;
        overflow: hidden;
      }

      #dropdown #header div {
        text-wrap: nowrap;
        overflow: hidden;
      }

      #dropdown #actions {
        margin-top: 0.5rem;
        display: flex;
        flex-direction: column;
        gap: 2px;
      }

      #dropdown #actions a,
      #dropdown #actions button {
        cursor: pointer;
        background-color: #fff;
        display: flex;
        border: 0;
        align-items: center;
        gap: 1rem;
        border-radius: 0.45rem;
        font-size: 0.85rem;
        padding: 0.8rem 0;
        text-align: center;
        text-wrap: nowrap;
        text-decoration: none;
        color: #000;
        transition: 200ms;
      }

      #dropdown #actions button:hover,
      #dropdown #actions button:focus-visible {
        background-color: #f8f8f8;
        padding: 0.8rem;
        transition: 200ms;
      }

      #dropdown #actions button:focus-visible {
        outline: 2px solid var(--primary-600);
      }

      #dropdown #actions .row {
        display: flex;
        justify-content: flex-end;
      }

      :host([location-vertical="top"]) #dropdown {
        top: auto;
        bottom: 100%;
        margin-top: 0;
        margin-bottom: 1rem;
        border-radius: 0.85rem 0 0 0;
        border-bottom: 2px solid var(--primary-700, #1448e1);
        border-top: 1px solid var(--gray-100, #eaeaea);
      }

      :host([location-vertical="top"][location-horizontal="left"]) #dropdown {
        border-radius: 0 0.85rem 0 0;
      }

      :host([location-vertical="bottom"][location-horizontal="left"])
        #dropdown {
        border-radius: 0 0 0.85rem 0;
      }

      :host([location-horizontal="left"]) #dropdown {
        left: 0;
      }

      :host([location-horizontal="right"]) #dropdown {
        left: calc(-1 * var(--width) - 1rem);
      }

      #launchpad {
        box-sizing: border-box;
        position: fixed;
        bottom: 0;
        left: 0;
        width: 100svw;
        padding: 1rem;
        background-color: #fff;
        color: black;
        font-size: 1rem;
        z-index: 1000;

        border-top: 2px solid var(--primary-800, #173ab6);

        display: flex;
        gap: 1rem;
        align-items: center;

        justify-content: space-between;
      }

      #launchpad div {
        display: flex;
        gap: 1rem;
        align-items: center;
      }

      #launchpad #launchpad-status {
        font-weight: 600;
      }

      #launchpad * {
        margin: 0;
        padding: 0;
      }

      #launchpad a {
        color: black;
      }

      #launchpad button {
        color: black;
        margin: 0;
        padding: 0;
        border: none;
        background-color: transparent;
        font-size: 1.5rem;
        cursor: pointer;
      }
    `],le([dt()],de.prototype,"location",void 0),le([ct()],de.prototype,"open",void 0),le([ct()],de.prototype,"launchpadForceClosed",void 0),de=le([at("locksmith-user-icon")],de);export{qt as AboutMeState,Ct as GenerateFingerprint,ae as LocksmithLayout,Yt as LocksmithLoginComponent,re as LocksmithProfileComponent,Gt as LocksmithRegistrationComponent,Qt as LocksmithResetPasswordComponent,de as LocksmithUserIconComponent,ne as SecureFetch,Kt as aboutMe,ht as inputStyles};
//# sourceMappingURL=locksmith-ui.bundle.js.map
