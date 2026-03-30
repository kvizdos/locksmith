/**
 * @license
 * Copyright 2019 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */
const e=globalThis,t=e.ShadowRoot&&(void 0===e.ShadyCSS||e.ShadyCSS.nativeShadow)&&"adoptedStyleSheets"in Document.prototype&&"replace"in CSSStyleSheet.prototype,i=Symbol(),o=new WeakMap;let r=class{constructor(e,t,o){if(this._$cssResult$=!0,o!==i)throw Error("CSSResult is not constructable. Use `unsafeCSS` or `css` instead.");this.cssText=e,this.t=t}get styleSheet(){let e=this.o;const i=this.t;if(t&&void 0===e){const t=void 0!==i&&1===i.length;t&&(e=o.get(i)),void 0===e&&((this.o=e=new CSSStyleSheet).replaceSync(this.cssText),t&&o.set(i,e))}return e}toString(){return this.cssText}};const s=(e,...t)=>{const o=1===e.length?e[0]:t.reduce(((t,i,o)=>t+(e=>{if(!0===e._$cssResult$)return e.cssText;if("number"==typeof e)return e;throw Error("Value passed to 'css' function must be a 'css' function result: "+e+". Use 'unsafeCSS' to pass non-literal values, but take care to ensure page security.")})(i)+e[o+1]),e[0]);return new r(o,e,i)},a=t?e=>e:e=>e instanceof CSSStyleSheet?(e=>{let t="";for(const i of e.cssRules)t+=i.cssText;return(e=>new r("string"==typeof e?e:e+"",void 0,i))(t)})(e):e
/**
 * @license
 * Copyright 2017 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */,{is:n,defineProperty:l,getOwnPropertyDescriptor:d,getOwnPropertyNames:c,getOwnPropertySymbols:h,getPrototypeOf:p}=Object,u=globalThis,v=u.trustedTypes,f=v?v.emptyScript:"",g=u.reactiveElementPolyfillSupport,m=(e,t)=>e,w={toAttribute(e,t){switch(t){case Boolean:e=e?f:null;break;case Object:case Array:e=null==e?e:JSON.stringify(e)}return e},fromAttribute(e,t){let i=e;switch(t){case Boolean:i=null!==e;break;case Number:i=null===e?null:Number(e);break;case Object:case Array:try{i=JSON.parse(e)}catch(e){i=null}}return i}},y=(e,t)=>!n(e,t),b={attribute:!0,type:String,converter:w,reflect:!1,hasChanged:y};Symbol.metadata??=Symbol("metadata"),u.litPropertyMetadata??=new WeakMap;class x extends HTMLElement{static addInitializer(e){this._$Ei(),(this.l??=[]).push(e)}static get observedAttributes(){return this.finalize(),this._$Eh&&[...this._$Eh.keys()]}static createProperty(e,t=b){if(t.state&&(t.attribute=!1),this._$Ei(),this.elementProperties.set(e,t),!t.noAccessor){const i=Symbol(),o=this.getPropertyDescriptor(e,i,t);void 0!==o&&l(this.prototype,e,o)}}static getPropertyDescriptor(e,t,i){const{get:o,set:r}=d(this.prototype,e)??{get(){return this[t]},set(e){this[t]=e}};return{get(){return o?.call(this)},set(t){const s=o?.call(this);r.call(this,t),this.requestUpdate(e,s,i)},configurable:!0,enumerable:!0}}static getPropertyOptions(e){return this.elementProperties.get(e)??b}static _$Ei(){if(this.hasOwnProperty(m("elementProperties")))return;const e=p(this);e.finalize(),void 0!==e.l&&(this.l=[...e.l]),this.elementProperties=new Map(e.elementProperties)}static finalize(){if(this.hasOwnProperty(m("finalized")))return;if(this.finalized=!0,this._$Ei(),this.hasOwnProperty(m("properties"))){const e=this.properties,t=[...c(e),...h(e)];for(const i of t)this.createProperty(i,e[i])}const e=this[Symbol.metadata];if(null!==e){const t=litPropertyMetadata.get(e);if(void 0!==t)for(const[e,i]of t)this.elementProperties.set(e,i)}this._$Eh=new Map;for(const[e,t]of this.elementProperties){const i=this._$Eu(e,t);void 0!==i&&this._$Eh.set(i,e)}this.elementStyles=this.finalizeStyles(this.styles)}static finalizeStyles(e){const t=[];if(Array.isArray(e)){const i=new Set(e.flat(1/0).reverse());for(const e of i)t.unshift(a(e))}else void 0!==e&&t.push(a(e));return t}static _$Eu(e,t){const i=t.attribute;return!1===i?void 0:"string"==typeof i?i:"string"==typeof e?e.toLowerCase():void 0}constructor(){super(),this._$Ep=void 0,this.isUpdatePending=!1,this.hasUpdated=!1,this._$Em=null,this._$Ev()}_$Ev(){this._$ES=new Promise((e=>this.enableUpdating=e)),this._$AL=new Map,this._$E_(),this.requestUpdate(),this.constructor.l?.forEach((e=>e(this)))}addController(e){(this._$EO??=new Set).add(e),void 0!==this.renderRoot&&this.isConnected&&e.hostConnected?.()}removeController(e){this._$EO?.delete(e)}_$E_(){const e=new Map,t=this.constructor.elementProperties;for(const i of t.keys())this.hasOwnProperty(i)&&(e.set(i,this[i]),delete this[i]);e.size>0&&(this._$Ep=e)}createRenderRoot(){const i=this.shadowRoot??this.attachShadow(this.constructor.shadowRootOptions);return((i,o)=>{if(t)i.adoptedStyleSheets=o.map((e=>e instanceof CSSStyleSheet?e:e.styleSheet));else for(const t of o){const o=document.createElement("style"),r=e.litNonce;void 0!==r&&o.setAttribute("nonce",r),o.textContent=t.cssText,i.appendChild(o)}})(i,this.constructor.elementStyles),i}connectedCallback(){this.renderRoot??=this.createRenderRoot(),this.enableUpdating(!0),this._$EO?.forEach((e=>e.hostConnected?.()))}enableUpdating(e){}disconnectedCallback(){this._$EO?.forEach((e=>e.hostDisconnected?.()))}attributeChangedCallback(e,t,i){this._$AK(e,i)}_$EC(e,t){const i=this.constructor.elementProperties.get(e),o=this.constructor._$Eu(e,i);if(void 0!==o&&!0===i.reflect){const r=(void 0!==i.converter?.toAttribute?i.converter:w).toAttribute(t,i.type);this._$Em=e,null==r?this.removeAttribute(o):this.setAttribute(o,r),this._$Em=null}}_$AK(e,t){const i=this.constructor,o=i._$Eh.get(e);if(void 0!==o&&this._$Em!==o){const e=i.getPropertyOptions(o),r="function"==typeof e.converter?{fromAttribute:e.converter}:void 0!==e.converter?.fromAttribute?e.converter:w;this._$Em=o,this[o]=r.fromAttribute(t,e.type),this._$Em=null}}requestUpdate(e,t,i){if(void 0!==e){if(i??=this.constructor.getPropertyOptions(e),!(i.hasChanged??y)(this[e],t))return;this.P(e,t,i)}!1===this.isUpdatePending&&(this._$ES=this._$ET())}P(e,t,i){this._$AL.has(e)||this._$AL.set(e,t),!0===i.reflect&&this._$Em!==e&&(this._$Ej??=new Set).add(e)}async _$ET(){this.isUpdatePending=!0;try{await this._$ES}catch(e){Promise.reject(e)}const e=this.scheduleUpdate();return null!=e&&await e,!this.isUpdatePending}scheduleUpdate(){return this.performUpdate()}performUpdate(){if(!this.isUpdatePending)return;if(!this.hasUpdated){if(this.renderRoot??=this.createRenderRoot(),this._$Ep){for(const[e,t]of this._$Ep)this[e]=t;this._$Ep=void 0}const e=this.constructor.elementProperties;if(e.size>0)for(const[t,i]of e)!0!==i.wrapped||this._$AL.has(t)||void 0===this[t]||this.P(t,this[t],i)}let e=!1;const t=this._$AL;try{e=this.shouldUpdate(t),e?(this.willUpdate(t),this._$EO?.forEach((e=>e.hostUpdate?.())),this.update(t)):this._$EU()}catch(t){throw e=!1,this._$EU(),t}e&&this._$AE(t)}willUpdate(e){}_$AE(e){this._$EO?.forEach((e=>e.hostUpdated?.())),this.hasUpdated||(this.hasUpdated=!0,this.firstUpdated(e)),this.updated(e)}_$EU(){this._$AL=new Map,this.isUpdatePending=!1}get updateComplete(){return this.getUpdateComplete()}getUpdateComplete(){return this._$ES}shouldUpdate(e){return!0}update(e){this._$Ej&&=this._$Ej.forEach((e=>this._$EC(e,this[e]))),this._$EU()}updated(e){}firstUpdated(e){}}x.elementStyles=[],x.shadowRootOptions={mode:"open"},x[m("elementProperties")]=new Map,x[m("finalized")]=new Map,g?.({ReactiveElement:x}),(u.reactiveElementVersions??=[]).push("2.0.4");
/**
 * @license
 * Copyright 2017 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */
const $=globalThis,k=$.trustedTypes,A=k?k.createPolicy("lit-html",{createHTML:e=>e}):void 0,_="$lit$",P=`lit$${Math.random().toFixed(9).slice(2)}$`,R="?"+P,S=`<${R}>`,M=document,O=()=>M.createComment(""),E=e=>null===e||"object"!=typeof e&&"function"!=typeof e,C=Array.isArray,z="[ \t\n\f\r]",L=/<(?:(!--|\/[^a-zA-Z])|(\/?[a-zA-Z][^>\s]*)|(\/?$))/g,T=/-->/g,U=/>/g,j=RegExp(`>|${z}(?:([^\\s"'>=/]+)(${z}*=${z}*(?:[^ \t\n\f\r"'\`<>=]|("|')|))|$)`,"g"),N=/'/g,H=/"/g,B=/^(?:script|style|textarea|title)$/i,D=(e=>(t,...i)=>({_$litType$:e,strings:t,values:i}))(1),I=Symbol.for("lit-noChange"),V=Symbol.for("lit-nothing"),Y=new WeakMap,q=M.createTreeWalker(M,129);function F(e,t){if(!C(e)||!e.hasOwnProperty("raw"))throw Error("invalid template strings array");return void 0!==A?A.createHTML(t):t}const W=(e,t)=>{const i=e.length-1,o=[];let r,s=2===t?"<svg>":3===t?"<math>":"",a=L;for(let t=0;t<i;t++){const i=e[t];let n,l,d=-1,c=0;for(;c<i.length&&(a.lastIndex=c,l=a.exec(i),null!==l);)c=a.lastIndex,a===L?"!--"===l[1]?a=T:void 0!==l[1]?a=U:void 0!==l[2]?(B.test(l[2])&&(r=RegExp("</"+l[2],"g")),a=j):void 0!==l[3]&&(a=j):a===j?">"===l[0]?(a=r??L,d=-1):void 0===l[1]?d=-2:(d=a.lastIndex-l[2].length,n=l[1],a=void 0===l[3]?j:'"'===l[3]?H:N):a===H||a===N?a=j:a===T||a===U?a=L:(a=j,r=void 0);const h=a===j&&e[t+1].startsWith("/>")?" ":"";s+=a===L?i+S:d>=0?(o.push(n),i.slice(0,d)+_+i.slice(d)+P+h):i+P+(-2===d?t:h)}return[F(e,s+(e[i]||"<?>")+(2===t?"</svg>":3===t?"</math>":"")),o]};class K{constructor({strings:e,_$litType$:t},i){let o;this.parts=[];let r=0,s=0;const a=e.length-1,n=this.parts,[l,d]=W(e,t);if(this.el=K.createElement(l,i),q.currentNode=this.el.content,2===t||3===t){const e=this.el.content.firstChild;e.replaceWith(...e.childNodes)}for(;null!==(o=q.nextNode())&&n.length<a;){if(1===o.nodeType){if(o.hasAttributes())for(const e of o.getAttributeNames())if(e.endsWith(_)){const t=d[s++],i=o.getAttribute(e).split(P),a=/([.?@])?(.*)/.exec(t);n.push({type:1,index:r,name:a[2],strings:i,ctor:"."===a[1]?X:"?"===a[1]?ee:"@"===a[1]?te:Q}),o.removeAttribute(e)}else e.startsWith(P)&&(n.push({type:6,index:r}),o.removeAttribute(e));if(B.test(o.tagName)){const e=o.textContent.split(P),t=e.length-1;if(t>0){o.textContent=k?k.emptyScript:"";for(let i=0;i<t;i++)o.append(e[i],O()),q.nextNode(),n.push({type:2,index:++r});o.append(e[t],O())}}}else if(8===o.nodeType)if(o.data===R)n.push({type:2,index:r});else{let e=-1;for(;-1!==(e=o.data.indexOf(P,e+1));)n.push({type:7,index:r}),e+=P.length-1}r++}}static createElement(e,t){const i=M.createElement("template");return i.innerHTML=e,i}}function J(e,t,i=e,o){if(t===I)return t;let r=void 0!==o?i._$Co?.[o]:i._$Cl;const s=E(t)?void 0:t._$litDirective$;return r?.constructor!==s&&(r?._$AO?.(!1),void 0===s?r=void 0:(r=new s(e),r._$AT(e,i,o)),void 0!==o?(i._$Co??=[])[o]=r:i._$Cl=r),void 0!==r&&(t=J(e,r._$AS(e,t.values),r,o)),t}class G{constructor(e,t){this._$AV=[],this._$AN=void 0,this._$AD=e,this._$AM=t}get parentNode(){return this._$AM.parentNode}get _$AU(){return this._$AM._$AU}u(e){const{el:{content:t},parts:i}=this._$AD,o=(e?.creationScope??M).importNode(t,!0);q.currentNode=o;let r=q.nextNode(),s=0,a=0,n=i[0];for(;void 0!==n;){if(s===n.index){let t;2===n.type?t=new Z(r,r.nextSibling,this,e):1===n.type?t=new n.ctor(r,n.name,n.strings,this,e):6===n.type&&(t=new ie(r,this,e)),this._$AV.push(t),n=i[++a]}s!==n?.index&&(r=q.nextNode(),s++)}return q.currentNode=M,o}p(e){let t=0;for(const i of this._$AV)void 0!==i&&(void 0!==i.strings?(i._$AI(e,i,t),t+=i.strings.length-2):i._$AI(e[t])),t++}}class Z{get _$AU(){return this._$AM?._$AU??this._$Cv}constructor(e,t,i,o){this.type=2,this._$AH=V,this._$AN=void 0,this._$AA=e,this._$AB=t,this._$AM=i,this.options=o,this._$Cv=o?.isConnected??!0}get parentNode(){let e=this._$AA.parentNode;const t=this._$AM;return void 0!==t&&11===e?.nodeType&&(e=t.parentNode),e}get startNode(){return this._$AA}get endNode(){return this._$AB}_$AI(e,t=this){e=J(this,e,t),E(e)?e===V||null==e||""===e?(this._$AH!==V&&this._$AR(),this._$AH=V):e!==this._$AH&&e!==I&&this._(e):void 0!==e._$litType$?this.$(e):void 0!==e.nodeType?this.T(e):(e=>C(e)||"function"==typeof e?.[Symbol.iterator])(e)?this.k(e):this._(e)}O(e){return this._$AA.parentNode.insertBefore(e,this._$AB)}T(e){this._$AH!==e&&(this._$AR(),this._$AH=this.O(e))}_(e){this._$AH!==V&&E(this._$AH)?this._$AA.nextSibling.data=e:this.T(M.createTextNode(e)),this._$AH=e}$(e){const{values:t,_$litType$:i}=e,o="number"==typeof i?this._$AC(e):(void 0===i.el&&(i.el=K.createElement(F(i.h,i.h[0]),this.options)),i);if(this._$AH?._$AD===o)this._$AH.p(t);else{const e=new G(o,this),i=e.u(this.options);e.p(t),this.T(i),this._$AH=e}}_$AC(e){let t=Y.get(e.strings);return void 0===t&&Y.set(e.strings,t=new K(e)),t}k(e){C(this._$AH)||(this._$AH=[],this._$AR());const t=this._$AH;let i,o=0;for(const r of e)o===t.length?t.push(i=new Z(this.O(O()),this.O(O()),this,this.options)):i=t[o],i._$AI(r),o++;o<t.length&&(this._$AR(i&&i._$AB.nextSibling,o),t.length=o)}_$AR(e=this._$AA.nextSibling,t){for(this._$AP?.(!1,!0,t);e&&e!==this._$AB;){const t=e.nextSibling;e.remove(),e=t}}setConnected(e){void 0===this._$AM&&(this._$Cv=e,this._$AP?.(e))}}class Q{get tagName(){return this.element.tagName}get _$AU(){return this._$AM._$AU}constructor(e,t,i,o,r){this.type=1,this._$AH=V,this._$AN=void 0,this.element=e,this.name=t,this._$AM=o,this.options=r,i.length>2||""!==i[0]||""!==i[1]?(this._$AH=Array(i.length-1).fill(new String),this.strings=i):this._$AH=V}_$AI(e,t=this,i,o){const r=this.strings;let s=!1;if(void 0===r)e=J(this,e,t,0),s=!E(e)||e!==this._$AH&&e!==I,s&&(this._$AH=e);else{const o=e;let a,n;for(e=r[0],a=0;a<r.length-1;a++)n=J(this,o[i+a],t,a),n===I&&(n=this._$AH[a]),s||=!E(n)||n!==this._$AH[a],n===V?e=V:e!==V&&(e+=(n??"")+r[a+1]),this._$AH[a]=n}s&&!o&&this.j(e)}j(e){e===V?this.element.removeAttribute(this.name):this.element.setAttribute(this.name,e??"")}}class X extends Q{constructor(){super(...arguments),this.type=3}j(e){this.element[this.name]=e===V?void 0:e}}class ee extends Q{constructor(){super(...arguments),this.type=4}j(e){this.element.toggleAttribute(this.name,!!e&&e!==V)}}class te extends Q{constructor(e,t,i,o,r){super(e,t,i,o,r),this.type=5}_$AI(e,t=this){if((e=J(this,e,t,0)??V)===I)return;const i=this._$AH,o=e===V&&i!==V||e.capture!==i.capture||e.once!==i.once||e.passive!==i.passive,r=e!==V&&(i===V||o);o&&this.element.removeEventListener(this.name,this,i),r&&this.element.addEventListener(this.name,this,e),this._$AH=e}handleEvent(e){"function"==typeof this._$AH?this._$AH.call(this.options?.host??this.element,e):this._$AH.handleEvent(e)}}class ie{constructor(e,t,i){this.element=e,this.type=6,this._$AN=void 0,this._$AM=t,this.options=i}get _$AU(){return this._$AM._$AU}_$AI(e){J(this,e)}}const oe=$.litHtmlPolyfillSupport;oe?.(K,Z),($.litHtmlVersions??=[]).push("3.2.1");
/**
 * @license
 * Copyright 2017 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */
let re=class extends x{constructor(){super(...arguments),this.renderOptions={host:this},this._$Do=void 0}createRenderRoot(){const e=super.createRenderRoot();return this.renderOptions.renderBefore??=e.firstChild,e}update(e){const t=this.render();this.hasUpdated||(this.renderOptions.isConnected=this.isConnected),super.update(e),this._$Do=((e,t,i)=>{const o=i?.renderBefore??t;let r=o._$litPart$;if(void 0===r){const e=i?.renderBefore??null;o._$litPart$=r=new Z(t.insertBefore(O(),e),e,void 0,i??{})}return r._$AI(e),r})(t,this.renderRoot,this.renderOptions)}connectedCallback(){super.connectedCallback(),this._$Do?.setConnected(!0)}disconnectedCallback(){super.disconnectedCallback(),this._$Do?.setConnected(!1)}render(){return I}};re._$litElement$=!0,re.finalized=!0,globalThis.litElementHydrateSupport?.({LitElement:re});const se=globalThis.litElementPolyfillSupport;se?.({LitElement:re}),(globalThis.litElementVersions??=[]).push("4.1.1");
/**
 * @license
 * Copyright 2017 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */
const ae=e=>(t,i)=>{void 0!==i?i.addInitializer((()=>{customElements.define(e,t)})):customElements.define(e,t)}
/**
 * @license
 * Copyright 2017 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */,ne={attribute:!0,type:String,converter:w,reflect:!1,hasChanged:y},le=(e=ne,t,i)=>{const{kind:o,metadata:r}=i;let s=globalThis.litPropertyMetadata.get(r);if(void 0===s&&globalThis.litPropertyMetadata.set(r,s=new Map),s.set(i.name,e),"accessor"===o){const{name:o}=i;return{set(i){const r=t.get.call(this);t.set.call(this,i),this.requestUpdate(o,r,e)},init(t){return void 0!==t&&this.P(o,void 0,e),t}}}if("setter"===o){const{name:o}=i;return function(i){const r=this[o];t.call(this,i),this.requestUpdate(o,r,e)}}throw Error("Unsupported decorator location: "+o)};function de(e){return(t,i)=>"object"==typeof i?le(e,t,i):((e,t,i)=>{const o=t.hasOwnProperty(i);return t.constructor.createProperty(i,o?{...e,wrapped:!0}:e),o?Object.getOwnPropertyDescriptor(t,i):void 0})(e,t,i)
/**
 * @license
 * Copyright 2017 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */}function ce(e){return de({...e,state:!0,attribute:!1})}const he=s`
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
 */,pe=1,ue=2,ve=e=>(...t)=>({_$litDirective$:e,values:t});class fe{constructor(e){}get _$AU(){return this._$AM._$AU}_$AT(e,t,i){this._$Ct=e,this._$AM=t,this._$Ci=i}_$AS(e,t){return this.update(e,t)}update(e,t){return this.render(...t)}}
/**
 * @license
 * Copyright 2018 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */const ge=ve(class extends fe{constructor(e){if(super(e),e.type!==pe||"class"!==e.name||e.strings?.length>2)throw Error("`classMap()` can only be used in the `class` attribute and must be the only part in the attribute.")}render(e){return" "+Object.keys(e).filter((t=>e[t])).join(" ")+" "}update(e,[t]){if(void 0===this.st){this.st=new Set,void 0!==e.strings&&(this.nt=new Set(e.strings.join(" ").split(/\s/).filter((e=>""!==e))));for(const e in t)t[e]&&!this.nt?.has(e)&&this.st.add(e);return this.render(t)}const i=e.element.classList;for(const e of this.st)e in t||(i.remove(e),this.st.delete(e));for(const e in t){const o=!!t[e];o===this.st.has(e)||this.nt?.has(e)||(o?(i.add(e),this.st.add(e)):(i.remove(e),this.st.delete(e)))}return I}});var me=function(e,t,i,o){var r,s=arguments.length,a=s<3?t:null===o?o=Object.getOwnPropertyDescriptor(t,i):o;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(e,t,i,o);else for(var n=e.length-1;n>=0;n--)(r=e[n])&&(a=(s<3?r(a):s>3?r(t,i,a):r(t,i))||a);return s>3&&a&&Object.defineProperty(t,i,a),a};let we=class extends re{constructor(){super(...arguments),this.disabled=!1,this.expectLoad=!1,this.loading=!1,this.loadingText=""}render(){return D`<button
      ?disabled=${this.disabled||this.loading}
      class=${ge({loading:this.loading})}
      @click=${()=>{this.dispatchEvent(new Event("fl-click"))}}
    >
      ${this.expectLoad?D`
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
          `:D``}
      ${this.loading&&""!==this.loadingText?this.loadingText:D`<slot></slot>`}
    </button>`}};we.styles=s`
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
  `,me([de()],we.prototype,"disabled",void 0),me([de()],we.prototype,"expectLoad",void 0),me([de()],we.prototype,"loading",void 0),me([de()],we.prototype,"loadingText",void 0),we=me([ae("button-component")],we);
/**
 * @license
 * Copyright 2020 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */
const ye=(e,t)=>{const i=e._$AN;if(void 0===i)return!1;for(const e of i)e._$AO?.(t,!1),ye(e,t);return!0},be=e=>{let t,i;do{if(void 0===(t=e._$AM))break;i=t._$AN,i.delete(e),e=t}while(0===i?.size)},xe=e=>{for(let t;t=e._$AM;e=t){let i=t._$AN;if(void 0===i)t._$AN=i=new Set;else if(i.has(e))break;i.add(e),Ae(t)}};
/**
 * @license
 * Copyright 2017 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */function $e(e){void 0!==this._$AN?(be(this),this._$AM=e,xe(this)):this._$AM=e}function ke(e,t=!1,i=0){const o=this._$AH,r=this._$AN;if(void 0!==r&&0!==r.size)if(t)if(Array.isArray(o))for(let e=i;e<o.length;e++)ye(o[e],!1),be(o[e]);else null!=o&&(ye(o,!1),be(o));else ye(this,e)}const Ae=e=>{e.type==ue&&(e._$AP??=ke,e._$AQ??=$e)};class _e extends fe{constructor(){super(...arguments),this._$AN=void 0}_$AT(e,t,i){super._$AT(e,t,i),xe(this),this.isConnected=e._$AU}_$AO(e,t=!0){e!==this.isConnected&&(this.isConnected=e,e?this.reconnected?.():this.disconnected?.()),t&&(ye(this,e),be(this))}setValue(e){if((e=>void 0===e.strings)(this._$Ct))this._$Ct._$AI(e,this);else{const t=[...this._$Ct._$AH];t[this._$Ci]=e,this._$Ct._$AI(t,this,0)}}disconnected(){}reconnected(){}}
/**
 * @license
 * Copyright 2020 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */const Pe=()=>new Re;class Re{}const Se=new WeakMap,Me=ve(class extends _e{render(e){return V}update(e,[t]){const i=t!==this.Y;return i&&void 0!==this.Y&&this.rt(void 0),(i||this.lt!==this.ct)&&(this.Y=t,this.ht=e.options?.host,this.rt(this.ct=e.element)),V}rt(e){if(this.isConnected||(e=void 0),"function"==typeof this.Y){const t=this.ht??globalThis;let i=Se.get(t);void 0===i&&(i=new WeakMap,Se.set(t,i)),void 0!==i.get(this.Y)&&this.Y.call(this.ht,void 0),i.set(this.Y,e),void 0!==e&&this.Y.call(this.ht,e)}else this.Y.value=e}get lt(){return"function"==typeof this.Y?Se.get(this.ht??globalThis)?.get(this.Y):this.Y?.value}disconnected(){this.lt===this.ct&&this.rt(void 0)}reconnected(){this.rt(this.ct)}});async function Oe(e){const t=(new TextEncoder).encode(e),i=await crypto.subtle.digest("SHA-256",t);return Array.from(new Uint8Array(i)).map((e=>e.toString(16).padStart(2,"0"))).join("")}async function Ee(){const e=window.screen.height*window.devicePixelRatio,t=window.screen.width*window.devicePixelRatio,i=window.screen.colorDepth;console.log("Color Depth",i);const o=await Oe(JSON.stringify({screenHeight:e,screenWidth:t,colorDepth:i})),r=Intl.DateTimeFormat().resolvedOptions().timeZone,s=window.navigator.hardwareConcurrency,a=window.navigator.language,n=await(async()=>{const e=document.createElement("canvas");e.width=500,e.height=500,e.style.display="none",document.body.appendChild(e);const t=e.getContext("2d");t.textBaseline="top",t.font="14px 'Arial'",t.textBaseline="alphabetic",t.fillStyle="#f60",t.fillRect(125,1,62,20),t.fillStyle="#069",t.fillText("Hello, world!",2,15),t.fillStyle="rgba(102, 204, 0, 0.7)",t.fillText("Hello, world!",4,17),t.fillText("🤙",100,20),t.fillText("🎉",110,25),t.fillText("🤣",115,30);const i=e.toDataURL();e.remove();return await Oe(i)})(),l=await(async()=>{const e=(()=>{const e=document.createElement("canvas");let t;try{t=e.getContext("webgl")||e.getContext("experimental-webgl")}catch(e){console.error("Failed to get WebGL context: ",e)}return t})();if(!e)return null;const t=e.getExtension("WEBGL_debug_renderer_info");if(t){const i={renderer:e.getParameter(t.UNMASKED_RENDERER_WEBGL),vendor:e.getParameter(t.UNMASKED_VENDOR_WEBGL)};return await Oe(JSON.stringify(i))}return await Oe(JSON.stringify("blank"))})(),d=await Oe(JSON.stringify({touchSupport:"ontouchstart"in window||navigator.maxTouchPoints>0,maxTouchPoints:navigator.maxTouchPoints})),c=(null===navigator||void 0===navigator?void 0:navigator.platform)||"unknown",h=await(async()=>{const e=new OfflineAudioContext(1,44100,44100),t=e.createOscillator();t.type="sine",t.frequency.setValueAtTime(1e3,e.currentTime),t.connect(e.destination),t.start(0);const i=(await e.startRendering()).getChannelData(0),o=new Uint8Array(i.length);for(let e=0;e<i.length;e++)o[e]=Math.floor(255*(.5*i[e]+.5));const r=await crypto.subtle.digest("SHA-256",o),s=Array.from(new Uint8Array(r)).map((e=>e.toString(16).padStart(2,"0"))).join("");return s})();let p={screen:o,timezone:r,hardwareConcurrency:s,deviceMemory:"0",canvas:n,lang:a,webgl:l,touch:d,battery:!1,platform:c,audio:h,userAgent:"",windowSize:null,dnt:null,devices:null};const u=await Oe(navigator.userAgent),v=await Oe(JSON.stringify({height:window.innerHeight,width:window.innerWidth})),f=navigator.doNotTrack||!1,g=await(async()=>{try{return(await navigator.mediaDevices.enumerateDevices()).map((e=>({kind:e.kind,label:e.label,deviceId:e.deviceId,groupId:e.groupId})))}catch(e){return console.log(e),[]}})(),m=await Oe(g.map((e=>Object.values(e).join(":"))).join(";"));return p={...p,userAgent:u,windowSize:v,dnt:f,devices:m},p}const Ce={activity:D`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 37 37"><path fill-rule="evenodd" d="M8 2h21a6 6 0 0 1 6 6v21a6 6 0 0 1-6 6H8a6 6 0 0 1-6-6V8a6 6 0 0 1 6-6M0 8a8 8 0 0 1 8-8h21a8 8 0 0 1 8 8v21a8 8 0 0 1-8 8H8a8 8 0 0 1-8-8zm8.5 1a1.5 1.5 0 1 0 0 3h21a1.5 1.5 0 0 0 0-3zM7 18.5A1.5 1.5 0 0 1 8.5 17h18a1.5 1.5 0 0 1 0 3h-18A1.5 1.5 0 0 1 7 18.5M8.5 25a1.5 1.5 0 0 0 0 3h20a1.5 1.5 0 0 0 0-3z" class="primary" clip-rule="evenodd"/></svg>`,alert:D`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 100 100"><g clip-path="url(#a)"><path stroke-width="5" d="m52.165 17.75 34.641 60c.962 1.667-.24 3.75-2.165 3.75H15.359c-1.924 0-3.127-2.083-2.165-3.75l34.64-60c.963-1.667 3.369-1.667 4.331 0" class="primary-stroke"/><path d="M44.414 40.384A5 5 0 0 1 49.4 35h1.202a5 5 0 0 1 4.985 5.383l-1.114 14.475a4.486 4.486 0 0 1-8.945 0z" class="primary"/><circle cx="50" cy="68" r="5" class="primary"/></g><defs><clipPath id="a"><path fill="#fff" d="M0 0h100v100H0z"/></clipPath></defs></svg>`,check:D`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 17 15"><path stroke-width="2" d="m1 9.5 3.695 3.695a1 1 0 0 0 1.5-.098L15.5 1" class="primary-stroke"/></svg>`,checkmark:D`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 33 33"><circle cx="16.15" cy="16.15" r="16.15" class="primary"/><path stroke="#fff" stroke-width="3" d="m8.604 18.867 3.328 3.328a1 1 0 0 0 1.452-.04L24.3 9.962"/></svg>`,clock:D`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 100 100"><path stroke-width="6" d="M68.5 14.526A39.8 39.8 0 0 0 50 10c-22.091 0-40 17.909-40 40s17.909 40 40 40 40-17.909 40-40c0-8.127-1.336-14.688-5.5-21" class="secondary-stroke"/><path d="M87.255 18.607a5 5 0 1 0-7.071-7.071L45.536 46.184a5 5 0 1 0 7.07 7.07zM24.16 82.33a5 5 0 0 0-8.66-5l-5 8.66a5 5 0 1 0 8.66 5zm51.34 0a5 5 0 1 1 8.66-5l5 8.66a5 5 0 0 1-8.66 5z" class="primary"/></svg>`,cog:D`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 95 95"><path fill-rule="evenodd" d="M43 0a5 5 0 0 0-5 5v8.286c0 1.856-1.237 3.473-2.951 4.185-1.715.712-3.71.432-5.024-.88l-5.86-5.86a5 5 0 0 0-7.07 0l-6.365 6.363a5 5 0 0 0 0 7.072l5.86 5.86c1.313 1.312 1.593 3.308.88 5.023C16.76 36.763 15.143 38 13.287 38H5a5 5 0 0 0-5 5v9a5 5 0 0 0 5 5h8.286c1.856 0 3.473 1.237 4.185 2.951.712 1.715.432 3.71-.88 5.024l-5.86 5.86a5 5 0 0 0 0 7.07l6.363 6.364a5 5 0 0 0 7.072 0l5.86-5.86c1.312-1.312 3.308-1.592 5.023-.88S38 79.858 38 81.714V90a5 5 0 0 0 5 5h9a5 5 0 0 0 5-5v-8.286c0-1.856 1.237-3.473 2.951-4.185 1.715-.712 3.71-.432 5.024.88l5.86 5.86a5 5 0 0 0 7.07 0l6.365-6.363a5 5 0 0 0 0-7.071l-5.86-5.86c-1.313-1.313-1.593-3.308-.88-5.024.71-1.714 2.327-2.951 4.183-2.951H90a5 5 0 0 0 5-5v-9a5 5 0 0 0-5-5h-8.286c-1.856 0-3.473-1.237-4.185-2.951-.712-1.715-.432-3.71.88-5.024l5.86-5.86a5 5 0 0 0 0-7.07l-6.363-6.365a5 5 0 0 0-7.071 0l-5.86 5.86c-1.313 1.313-3.308 1.593-5.024.88C58.237 16.76 57 15.143 57 13.287V5a5 5 0 0 0-5-5zm4 62c8.284 0 15-6.716 15-15s-6.716-15-15-15-15 6.716-15 15 6.716 15 15 15" class="primary" clip-rule="evenodd"/></svg>`,email:D`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 100 100"><path d="M48.19 50.952 7 30a6 6 0 0 1 6-6h74a6 6 0 0 1 6 6L52.765 50.93a5 5 0 0 1-4.574.022" class="primary"/><path fill-rule="evenodd" d="M88 26H12a4 4 0 0 0-4 4v41a4 4 0 0 0 4 4h76a4 4 0 0 0 4-4V30a4 4 0 0 0-4-4m-76-4a8 8 0 0 0-8 8v41a8 8 0 0 0 8 8h76a8 8 0 0 0 8-8V30a8 8 0 0 0-8-8z" class="secondary" clip-rule="evenodd"/></svg>`,flag:D`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 9 13"><path d="M8 5V1l-1.175.294a10 10 0 0 1-5.588-.215L1 1v4l.237.08a10 10 0 0 0 5.588.214z" class="secondary"/><path d="M1 12.5V5m0 0V1l.237.08a10 10 0 0 0 5.588.214L8 1v4l-1.175.294a10 10 0 0 1-5.588-.215z" class="primary-stroke"/></svg>`,home:D`<svg xmlns="http://www.w3.org/2000/svg" class="icon-home" viewBox="0 0 24 24"><path d="M9 22H5a1 1 0 0 1-1-1V11l8-8 8 8v10a1 1 0 0 1-1 1h-4a1 1 0 0 1-1-1v-4a1 1 0 0 0-1-1h-2a1 1 0 0 0-1 1v4a1 1 0 0 1-1 1m3-9a2 2 0 1 0 0-4 2 2 0 0 0 0 4" class="primary"/><path d="m12.01 4.42-8.3 8.3a1 1 0 1 1-1.42-1.41l9.02-9.02a1 1 0 0 1 1.41 0l8.99 9.02a1 1 0 0 1-1.42 1.41z" class="secondary"/></svg>`,info:D`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"><circle cx="12" cy="12" r="11.5" stroke="#fff"/><path stroke-width="2" d="M13.5 18.5V13a1 1 0 0 0-1-1H10m3.5 6.5h-4m4 0h3" class="primary-stroke"/><circle cx="12.5" cy="7" r="2" class="primary"/></svg>`,note:D`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 100 100"><g stroke-width="6" clip-path="url(#a)"><path d="M58.657 3H18C9.716 3 3 9.716 3 18v64c0 8.284 6.716 15 15 15h64c8.284 0 15-6.716 15-15V34.629" class="primary-stroke"/><path d="M48.93 54.861 79.801 3.473a1 1 0 0 1 1.358-.35L92.707 9.79a1 1 0 0 1 .406 1.29l-.049.091L62.38 62.25 42.86 76.275z" class="secondary-stroke"/></g><defs><clipPath id="a"><path fill="#fff" d="M0 0h100v100H0z"/></clipPath></defs></svg>`,"person-group":D`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M12 13a3 3 0 0 1 3-3h4a3 3 0 0 1 3 3v3a1 1 0 0 1-1 1h-8a1 1 0 0 1-1-1 1 1 0 0 1-1 1H3a1 1 0 0 1-1-1v-3a3 3 0 0 1 3-3h4a3 3 0 0 1 3 3M7 9a3 3 0 1 1 0-6 3 3 0 0 1 0 6m10 0a3 3 0 1 1 0-6 3 3 0 0 1 0 6" class="secondary"/><path d="M12 13a3 3 0 1 1 0-6 3 3 0 0 1 0 6m-3 1h6a3 3 0 0 1 3 3v3a1 1 0 0 1-1 1H7a1 1 0 0 1-1-1v-3a3 3 0 0 1 3-3" class="primary"/></svg>`,"person-outline":D`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 30 34"><path fill-rule="evenodd" d="M20 7a5 5 0 1 1-10 0 5 5 0 0 1 10 0m2 0A7 7 0 1 1 8 7a7 7 0 0 1 14 0M.187 21.912C-.438 17.746 2.788 14 7 14l1.694 1.376a10 10 0 0 0 12.612 0L23 14c4.212 0 7.438 3.746 6.813 7.912L28 34H2zm22.38-4.983 1.09-.886a4.89 4.89 0 0 1 4.178 5.572L26.278 32H3.722L2.165 21.615a4.89 4.89 0 0 1 4.178-5.572l1.09.886a12 12 0 0 0 15.134 0" class="primary" clip-rule="evenodd"/></svg>`,person:D`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 30 34"><g clip-path="url(#a)"><path fill-rule="evenodd" d="M20 7a5 5 0 1 1-10 0 5 5 0 0 1 10 0m2 0A7 7 0 1 1 8 7a7 7 0 0 1 14 0M.187 21.912C-.438 17.746 2.788 14 7 14l1.694 1.376a10 10 0 0 0 12.612 0L23 14c4.212 0 7.438 3.746 6.813 7.912L28 34H2z" class="primary" clip-rule="evenodd"/></g><defs><clipPath id="a"><path fill="#fff" d="M0 0h30v34H0z"/></clipPath></defs></svg>`,"phone-disabled":D`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 100 100"><g clip-path="url(#a)"><rect width="31" height="11.499" x="37.69" y="4.483" stroke="#A4A4A4" stroke-width="4" rx="3" transform="rotate(60 37.69 4.483)"/><path stroke="#A4A4A4" stroke-width="4" d="M56.212 88a13 13 0 0 1-9.46-4.082l-.233-.255-14.483-16.209a13 13 0 0 1-2.514-4.191l-.109-.31L20.205 35.7a13 13 0 0 1 1.186-10.876l.196-.315 3.355-5.218 9.737 16.23c.21.348.345.735.4 1.136l.018.174.88 11.26a27 27 0 0 0 12.767 20.893l.719.426 6.43 3.689c.383.22.713.52.965.88l.103.158L65.434 88z"/><rect width="31" height="11.499" x="70.69" y="60.732" stroke="#A4A4A4" stroke-width="4" rx="3" transform="rotate(60 70.69 60.732)"/><circle cx="36.869" cy="60.869" r="19.869" class="primary"/><path fill="#fff" fill-rule="evenodd" d="M30.32 49.486a1 1 0 0 0-1.413 0l-3.683 3.682a1 1 0 0 0 0 1.415l5.908 5.907a1 1 0 0 1 0 1.414l-6.103 6.103a1 1 0 0 0 0 1.414l3.55 3.55a1 1 0 0 0 1.414 0l6.103-6.103a1 1 0 0 1 1.414 0l5.907 5.908a1 1 0 0 0 1.415 0l3.682-3.682a1 1 0 0 0 0-1.415l-5.908-5.907a1 1 0 0 1 0-1.414l6.103-6.103a1 1 0 0 0 0-1.415l-3.55-3.55a1 1 0 0 0-1.413 0l-6.104 6.104a1 1 0 0 1-1.414 0z" clip-rule="evenodd"/></g><defs><clipPath id="a"><path fill="#fff" d="M0 0h100v100H0z"/></clipPath></defs></svg>`,phone:D`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 100 100"><g clip-path="url(#a)"><rect width="35" height="15.499" x="37.422" y="-.249" class="secondary" rx="5" transform="rotate(60 37.422 -.25)"/><path d="M24.13 16.854a1 1 0 0 1 1.698.026l10.566 17.61c.399.664.637 1.412.698 2.184l.88 11.26a25 25 0 0 0 12.486 19.74l6.431 3.689a5 5 0 0 1 1.779 1.73l9.402 15.386A1 1 0 0 1 67.217 90H56.212a15 15 0 0 1-11.185-5.005L30.544 68.787a15 15 0 0 1-3.026-5.193L18.311 36.34a15 15 0 0 1 1.593-12.913z" class="primary"/><rect width="35" height="15.499" x="70.422" y="56" class="secondary" rx="5" transform="rotate(60 70.422 56)"/></g><defs><clipPath id="a"><path fill="#fff" d="M0 0h100v100H0z"/></clipPath></defs></svg>`,pin:D`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 100 100"><path fill="#fff" d="M0 0h100v100H0z"/><path fill-rule="evenodd" d="M34.825 12.39A5 5 0 0 1 39.787 8H60.94a5 5 0 0 1 4.971 4.465l5.939 55.142a5 5 0 0 1-4.971 5.535h-5.264q.036-.456.036-.923c0-2.683-.914-5.153-2.447-7.116A5 5 0 0 0 62.37 59.9l-2.89-26.696a5 5 0 0 0-4.971-4.462H46.4a5 5 0 0 0-4.963 4.386l-3.302 26.697A5 5 0 0 0 41.045 65a11.52 11.52 0 0 0-2.493 8.142h-5.551a5 5 0 0 1-4.963-5.61z" class="primary" clip-rule="evenodd"/><circle cx="49.868" cy="72" r="7" class="secondary"/><rect width="8" height="18" x="46" y="75" class="secondary" rx="3"/></svg>`,search:D`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"><path fill-rule="evenodd" d="M9.92 7.93a5.93 5.93 0 1 1 11.858 0 5.93 5.93 0 0 1-11.859 0M15.848 0a7.93 7.93 0 0 0-6.27 12.785L.293 22.07l1.414 1.414 9.286-9.286A7.93 7.93 0 1 0 15.848 0" class="primary" clip-rule="evenodd"/></svg>`,"sign-out":D`<svg xmlns="http://www.w3.org/2000/svg" class="icon-door-exit" viewBox="0 0 24 24"><path d="M11 4h3a1 1 0 0 1 1 1v3a1 1 0 0 1-2 0V6h-2v12h2v-2a1 1 0 0 1 2 0v3a1 1 0 0 1-1 1h-3v1a1 1 0 0 1-1.27.96l-6.98-2A1 1 0 0 1 2 19V5a1 1 0 0 1 .75-.97l6.98-2A1 1 0 0 1 11 3z" class="primary"/><path d="m18.59 11-1.3-1.3c-.94-.94.47-2.35 1.42-1.4l3 3a1 1 0 0 1 0 1.4l-3 3c-.95.95-2.36-.46-1.42-1.4l1.3-1.3H14a1 1 0 0 1 0-2z" class="secondary"/></svg>`,sort:D`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 76 76"><rect width="70" height="11" x="3" y="16" class="primary" rx="5"/><rect width="62" height="11" x="11" y="33" class="primary" rx="5"/><rect width="54" height="11" x="19" y="50" class="primary" rx="5"/></svg>`,trash:D`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 57 58"><path stroke-width="2" d="M6.13 18.658h44.74a4 4 0 0 1 3.918 4.804l-6.023 29.356a4 4 0 0 1-4.232 3.184L28.97 54.778a6 6 0 0 0-.94 0l-15.563 1.224a4 4 0 0 1-4.232-3.184L2.212 23.462a4 4 0 0 1 3.918-4.805" class="primary-stroke"/><rect width="4" height="30.964" class="secondary" rx="2" transform="matrix(.99209 -.12553 .2006 .97967 9.295 22.952)"/><rect width="4" height="30.964" class="secondary" rx="2" transform="matrix(.9921 .12548 -.20051 .9797 44.157 22.45)"/><rect width="4" height="28.805" x="26.872" y="22.138" class="secondary" rx="2"/><path fill-rule="evenodd" d="M37.036 0a3.68 3.68 0 0 1 3.678 3.679 3.68 3.68 0 0 0 3.679 3.678h9.664a2.943 2.943 0 0 1 0 5.886H2.943a2.943 2.943 0 0 1 0-5.886h9.664a3.68 3.68 0 0 0 3.679-3.678A3.68 3.68 0 0 1 19.964 0zM22.564 2.207a2.207 2.207 0 1 0 0 4.415h11.872a2.207 2.207 0 0 0 0-4.415z" class="primary" clip-rule="evenodd"/></svg>`,unlink:D`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 100 100"><path stroke-width="8" d="m41.035 46-24.49 24.452a5 5 0 0 0 0 7.077l5.957 5.945a5 5 0 0 0 7.065 0L46 67.067M58.195 54l25.103-24.393a5 5 0 0 0-.01-7.183l-6.276-6.06a5 5 0 0 0-6.957.009L53 32.933" class="primary-stroke"/><rect width="8" height="18" x="65" y="74.997" class="shadow" rx="4" transform="rotate(-45 65 74.997)"/><rect width="8" height="18" x="73.498" y="63.489" class="shadow" rx="4" transform="rotate(-75 73.498 63.489)"/><rect width="8" height="18" x="49.681" y="79.357" class="shadow" rx="4" transform="rotate(-15 49.68 79.357)"/><rect width="8" height="18" x="34.445" y="21.543" class="shadow" rx="4" transform="rotate(135 34.445 21.543)"/><rect width="8" height="18" x="24.947" y="33.05" class="shadow" rx="4" transform="rotate(105 24.947 33.05)"/><rect width="8" height="18" x="49.765" y="18.182" class="shadow" rx="4" transform="rotate(165 49.765 18.182)"/></svg>`,"view-hidden":D`<svg xmlns="http://www.w3.org/2000/svg" class="icon-view-hidden" viewBox="0 0 24 24"><path d="M15.1 19.34a8 8 0 0 1-8.86-1.68L1.3 12.7a1 1 0 0 1 0-1.42L4.18 8.4l2.8 2.8a5 5 0 0 0 5.73 5.73l2.4 2.4zM8.84 4.6a8 8 0 0 1 8.7 1.74l4.96 4.95a1 1 0 0 1 0 1.42l-2.78 2.78-2.87-2.87a5 5 0 0 0-5.58-5.58L8.85 4.6z" class="primary"/><path d="m3.3 4.7 16 16a1 1 0 0 0 1.4-1.4l-16-16a1 1 0 0 0-1.4 1.4" class="secondary"/></svg>`,"view-visible":D`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M17.56 17.66a8 8 0 0 1-11.32 0L1.3 12.7a1 1 0 0 1 0-1.42l4.95-4.95a8 8 0 0 1 11.32 0l4.95 4.95a1 1 0 0 1 0 1.42l-4.95 4.95zM11.9 17a5 5 0 1 0 0-10 5 5 0 0 0 0 10" class="primary"/><circle cx="12" cy="12" r="3" class="secondary"/></svg>`,xmark:D`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 33 33"><circle cx="16.149" cy="16.149" r="16.149" class="primary"/><path stroke="#fff" stroke-width="3" d="m9.81 9.96 6.34 6.34m6.339 6.339-6.34-6.339m0 0 6.34-6.34m-6.34 6.34-6.338 6.339"/></svg>`};var ze,Le=function(e,t,i,o){var r,s=arguments.length,a=s<3?t:null===o?o=Object.getOwnPropertyDescriptor(t,i):o;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(e,t,i,o);else for(var n=e.length-1;n>=0;n--)(r=e[n])&&(a=(s<3?r(a):s>3?r(t,i,a):r(t,i))||a);return s>3&&a&&Object.defineProperty(t,i,a),a};let Te=ze=class extends re{constructor(){super(...arguments),this.name="info",this.size="24px",this.hoverable=!1,this.colorway="primary"}render(){var e;const t=null!==(e=ze.colorways[this.colorway])&&void 0!==e?e:ze.colorways.primary;return D`
      <div
        class=${ge({hoverable:this.hoverable})}
        style="
          --size: ${this.size};
          --primary: ${t.primary};
          --secondary: ${t.secondary};
          --shadow: ${t.shadow};
        "
      >
        ${Ce[this.name]}
      </div>
    `}};function Ue(e){return"function"==typeof e?e():e}Te.colorways={primary:{primary:"var(--primary-600)",secondary:"var(--primary-500, #327eff)",shadow:"var(--gray-400, #989898)"},danger:{primary:"var(--danger-600, red)",secondary:"var(--danger-500, pink)",shadow:"var(--gray-500, #888)"}},Te.styles=s`
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
  `,Le([de()],Te.prototype,"name",void 0),Le([de()],Te.prototype,"size",void 0),Le([de({type:Boolean})],Te.prototype,"hoverable",void 0),Le([de()],Te.prototype,"colorway",void 0),Te=ze=Le([ae("ui-icon")],Te);class je extends Event{static{this.eventName="lit-state-changed"}constructor(e,t,i){super(je.eventName,{cancelable:!1}),this.key=e,this.value=t,this.state=i}}const Ne=(e,t)=>t!==e&&(t==t||e==e);class He extends EventTarget{static{this.finalized=!1}static initPropertyMap(){this.propertyMap||(this.propertyMap=new Map)}get propertyMap(){return this.constructor.propertyMap}get stateValue(){return Object.fromEntries([...this.propertyMap].map((([e])=>[e,this[e]])))}constructor(){super(),this.hookMap=new Map,this.constructor.finalize(),this.propertyMap&&[...this.propertyMap].forEach((([e,t])=>{if(void 0!==t.initialValue){const i=Ue(t.initialValue);this[e]=i,t.value=i}}))}static finalize(){if(this.finalized)return!1;this.finalized=!0;const e=Object.keys(this.properties||{});for(const t of e)this.createProperty(t,this.properties[t]);return!0}static createProperty(e,t){this.finalize();const i="symbol"==typeof e?Symbol():`__${e}`,o=this.getPropertyDescriptor(String(e),i,t);Object.defineProperty(this.prototype,e,o)}static getPropertyDescriptor(e,t,i){const o=i?.hasChanged||Ne;return{get(){return this[t]},set(i){const r=this[e];this[t]=i,!0===o(i,r)&&this.dispatchStateEvent(e,i,this)},configurable:!0,enumerable:!0}}reset(){this.hookMap.forEach((e=>e.reset())),[...this.propertyMap].filter((([e,t])=>!(!0===t.skipReset||void 0===t.resetValue))).forEach((([e,t])=>{this[e]=t.resetValue}))}subscribe(e,t,i){t&&!Array.isArray(t)&&(t=[t]);const o=i=>{t&&!t.includes(i.key)||e(i.key,i.value,this)};return this.addEventListener(je.eventName,o,i),()=>this.removeEventListener(je.eventName,o)}dispatchStateEvent(e,t,i){this.dispatchEvent(new je(e,t,i))}}class Be{constructor(e,t,i){this.host=e,this.state=t,this.callback=i||(()=>this.host.requestUpdate()),this.host.addController(this)}hostConnected(){this.state.addEventListener(je.eventName,this.callback),this.callback()}hostDisconnected(){this.state.removeEventListener(je.eventName,this.callback)}}function De(e){return(t,i)=>{if(Object.getOwnPropertyDescriptor(t,i))throw new Error("@property must be called before all state decorators");const o=t.constructor;o.initPropertyMap();const r=t.hasOwnProperty(i);return o.propertyMap.set(i,{...e,initialValue:e?.value,resetValue:e?.value}),o.createProperty(i,e),r?Object.getOwnPropertyDescriptor(t,i):void 0}}new URL(location.href);const Ie={prefix:"_ls"};function Ve(e){return e={...Ie,...e},(t,i)=>{const o=Object.getOwnPropertyDescriptor(t,i);if(!o)throw new Error("@local-storage decorator need to be called after @property");const r=`${e?.prefix||""}_${e?.key||String(i)}`,s=t.constructor,a=s.propertyMap.get(i),n=a?.type;if(a){const t=a.initialValue;a.initialValue=()=>function(e,t){if(null!==e&&(t===Boolean||t===Number||t===Array||t===Object))try{e=JSON.parse(e)}catch(t){console.warn("cannot parse value",e)}return e}(localStorage.getItem(r),n)??Ue(t),s.propertyMap.set(i,{...a,...e})}const l=o?.set,d={...o,set:function(e){void 0!==e&&localStorage.setItem(r,n===Object||n===Array?JSON.stringify(e):e),l&&l.call(this,e)}};Object.defineProperty(s.prototype,i,d)}}var Ye=function(e,t,i,o){var r,s=arguments.length,a=s<3?t:null===o?o=Object.getOwnPropertyDescriptor(t,i):o;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(e,t,i,o);else for(var n=e.length-1;n>=0;n--)(r=e[n])&&(a=(s<3?r(a):s>3?r(t,i,a):r(t,i))||a);return s>3&&a&&Object.defineProperty(t,i,a),a};class qe extends He{get info(){return this.infoRaw?JSON.parse(this.infoRaw):void 0}set info(e){this.infoRaw=e?JSON.stringify(e):""}get username(){var e,t;return null!==(t=null===(e=this.info)||void 0===e?void 0:e.username.split("@")[0])&&void 0!==t?t:"unknown"}get email(){var e,t;return null!==(t=null===(e=this.info)||void 0===e?void 0:e.email)&&void 0!==t?t:"unknown"}get role(){var e,t;return null!==(t=null===(e=this.info)||void 0===e?void 0:e.role)&&void 0!==t?t:"unknown"}set extras(e){this.rawExtras=JSON.stringify(e)}get extras(){if(""===this.rawExtras||void 0===this.rawExtras)return{};try{return JSON.parse(this.rawExtras)}catch(e){return console.error(e),{}}}sendAboutMeUp(){this.aboutMeSent||(this.aboutMeSent=!0,window.dispatchEvent(new CustomEvent("locksmith-aboutme",{detail:{id:this.info.id,username:this.info.username,role:this.info.role,permissions:this.permissions.split(",")},bubbles:!0,composed:!0})))}async loadIfNeeded(e=!1){var t;const i=Date.now(),o=""!==this.infoRaw,r=(null===(t=this.permissions)||void 0===t?void 0:t.length)>0,s=0!=+this.loadedAt&&i-+this.loadedAt>9e5;if(!o||!r||s||e)try{const e=await fetch("/api/management/me");if(!e.ok)throw new Error("Failed to fetch /me");const t=await e.json();if(this.info=t.info,this.permissions=t.permissions.join(","),this.loadedAt=`${i}`,qe.GetExtraAboutMe){const e=await qe.GetExtraAboutMe({id:t.info.id,username:t.info.username,role:t.info.role,permissions:t.permissions});this.extras=e}this.sendAboutMeUp()}catch(e){console.error("Failed to load /me:",e),window.location.href=`/login?b=${encodeURIComponent(window.location.pathname+window.location.search)}&utm_source=locksmith&utm_campaign=session_expired`}else this.sendAboutMeUp()}hasPermission(e){return(this.permissions||"").split(",").includes(e)}hasRole(e){return Array.isArray(e)?e.some((e=>this.role===e)):this.role===e}get isLaunchpad(){const e=document.cookie.split(";");for(let t of e)if(t=t.trim(),t.startsWith("LaunchpadUser="))return t.substring(14)}clear(){localStorage.removeItem("_identity_i"),localStorage.removeItem("_identity_p"),localStorage.removeItem("_identity_la"),localStorage.removeItem("_identity_x")}signOut(){void 0!==qe.SignOutCallback&&qe.SignOutCallback({id:this.info.id,username:this.info.username,role:this.info.role,permissions:this.permissions.split(",")}),this.clear(),window.location.href="/sign-out"}}qe.SignOutCallback=void 0,qe.GetExtraAboutMe=void 0,Ye([Ve({key:"i",prefix:"_identity"}),De()],qe.prototype,"infoRaw",void 0),Ye([Ve({key:"x",prefix:"_identity"}),De()],qe.prototype,"rawExtras",void 0),Ye([Ve({key:"p",prefix:"_identity"}),De()],qe.prototype,"permissions",void 0),Ye([Ve({key:"la",prefix:"_identity"}),De()],qe.prototype,"loadedAt",void 0),Ye([De()],qe.prototype,"aboutMeSent",void 0);const Fe=new qe;var We=function(e,t,i,o){var r,s=arguments.length,a=s<3?t:null===o?o=Object.getOwnPropertyDescriptor(t,i):o;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(e,t,i,o);else for(var n=e.length-1;n>=0;n--)(r=e[n])&&(a=(s<3?r(a):s>3?r(t,i,a):r(t,i))||a);return s>3&&a&&Object.defineProperty(t,i,a),a};let Ke=class extends re{constructor(){super(...arguments),this.jsonSettings="",this.appName="FILL ME",this.settings={OauthProviders:[]},this.originOverride="",this.showingPassword=!1,this.errorMsg=void 0,this.loadedOnce=!1,this.isOnboarding=!1,this.loadingProvider=null,this.signInRef=Pe(),this.emailRef=Pe(),this.passwordRef=Pe(),this.aboutMeState=new Be(this,Fe)}connectedCallback(){super.connectedCallback();const e=new URLSearchParams(window.location.search);"true"===e.get("onboard")&&(this.isOnboarding=!0);const t=e.get("err");if(t&&"oauth_email_not_found"===t)this.errorMsg="No account is associated with that email address."}firstUpdated(){var e;null===(e=this.emailRef.value)||void 0===e||e.focus(),""!==this.jsonSettings&&setTimeout((()=>{try{this.settings=JSON.parse(this.jsonSettings)}catch(e){throw console.error("Invalid JSON in jsonSettings:",this.jsonSettings),e}this.loadedOnce=!0}))}canSignIn(){var e,t;return(null===(e=this.emailRef.value)||void 0===e?void 0:e.value.length)>0&&(null===(t=this.passwordRef.value)||void 0===t?void 0:t.value.length)>0}lockInputs(){this.emailRef.value.disabled=!0,this.passwordRef.value.disabled=!0}unlockInputs(){this.emailRef.value.disabled=!1,this.passwordRef.value.disabled=!1}async sendLoginRequest(){var e;const t=await Ee(),i={username:this.emailRef.value.value,password:this.passwordRef.value.value,fingerprint:t},o=await fetch(`${null!==(e=this.originOverride)&&void 0!==e?e:""}/api/login`,{method:"POST",body:JSON.stringify(i)});if(200!==o.status){const e=await o.json();if(e.error)throw new Error(e.error);throw new Error("Something went wrong.")}await Fe.loadIfNeeded(!0)}async attemptSignIn(){var e,t,i;if(this.errorMsg=void 0,!this.canSignIn())return this.errorMsg="Please enter your username and password.",void(0===(null===(e=this.emailRef.value)||void 0===e?void 0:e.value.length)?null===(t=this.emailRef.value)||void 0===t||t.focus():null===(i=this.passwordRef.value)||void 0===i||i.focus());this.signInRef.value.loading=!0,this.lockInputs(),this.requestUpdate();try{await this.sendLoginRequest();let e=this.isOnboarding&&void 0!==this.settings.PathToOnboard?this.settings.PathToOnboard:"/app";const t=new URLSearchParams(window.location.search).get("b");if(t){const i=decodeURIComponent(t);i.length>0&&"/"===i[0]&&(e=i)}window.location.href=e}catch(e){console.error(e),this.errorMsg=void 0,await this.updateComplete,this.errorMsg=e.message}finally{this.signInRef.value.loading=!1,this.unlockInputs()}}keydownEvent(e){var t,i;"Enter"===e.key&&(null===(t=this.emailRef.value)||void 0===t?void 0:t.value.length)>0&&(null===(i=this.passwordRef.value)||void 0===i?void 0:i.value.length)>0&&this.attemptSignIn()}render(){return D` <div id="root" class="${this.loadedOnce?"":"hide"}">
      <div id="header">
        <h1>Sign in to ${this.appName}</h1>
        ${!0!==this.settings.PublicRegistrationsDisabled||this.isOnboarding?D`
              <p id="intro">
                ${this.isOnboarding?D`<strong>Thank you for registering.</strong> Please sign
                      in for the first time to ensure everything works.`:D`Need an account?
                      <a href="/register">Create account</a>`}
              </p>
            `:D``}
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
            ${Me(this.emailRef)}
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
              @click=${()=>{const e=this.passwordRef.value.type;this.passwordRef.value.type="password"===e?"text":"password",this.showingPassword="text"!==e}}
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
            ${Me(this.passwordRef)}
            autofill="password"
            type="password"
            placeholder="Your password"
          />
        </div>
      </div>

      <div id="signInArea">
        <button-component
          ${Me(this.signInRef)}
          class="big"
          .expectLoad=${!0}
          .loadingText=${"Signing in.."}
          @fl-click=${this.attemptSignIn}
          >Sign in</button-component
        >

        <a href="/reset-password">Forgot Password</a>
      </div>

      ${this.settings.OauthProviders.length>0?D`
            <div class="hr-split">
              <hr />
              <p>Or...</p>
              <hr />
            </div>
          `:void 0}
      ${this.settings.OauthProviders.map((e=>D`
      <a class="oauth" href="/api/auth/oauth/${e}"
          @click=${t=>{t.preventDefault(),this.loadingProvider=e,window.location.href=t.currentTarget.href}}>
        <img src="/api/auth/oauth/${e}/logo"></img>
            ${this.loadingProvider===e?"Logging in...":`Sign in with ${e.charAt(0).toUpperCase()+e.slice(1)}`}
          <span></span></a>
      `))}
    </div>`}};Ke.styles=[he,s`
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
    `],We([de()],Ke.prototype,"jsonSettings",void 0),We([de()],Ke.prototype,"appName",void 0),We([de()],Ke.prototype,"settings",void 0),We([de()],Ke.prototype,"originOverride",void 0),We([ce()],Ke.prototype,"showingPassword",void 0),We([ce()],Ke.prototype,"errorMsg",void 0),We([ce()],Ke.prototype,"loadedOnce",void 0),We([ce()],Ke.prototype,"isOnboarding",void 0),We([ce()],Ke.prototype,"loadingProvider",void 0),Ke=We([ae("locksmith-login")],Ke);var Je=function(e,t,i,o){var r,s=arguments.length,a=s<3?t:null===o?o=Object.getOwnPropertyDescriptor(t,i):o;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(e,t,i,o);else for(var n=e.length-1;n>=0;n--)(r=e[n])&&(a=(s<3?r(a):s>3?r(t,i,a):r(t,i))||a);return s>3&&a&&Object.defineProperty(t,i,a),a};let Ge=class extends re{constructor(){super(...arguments),this.errorTitle="There was a bump in the road.",this.errorDescription="That's all we know."}connectedCallback(){super.connectedCallback()}render(){return D`<div>
      <div id="header">
        <ui-icon name="xmark" size="2rem"></ui-icon>
        <h1>${this.errorTitle}</h1>
        <p id="intro">${this.errorDescription}</p>
      </div>

      <hr />

      <button-component
        class="big"
        @fl-click=${()=>{window.location.href="/app"}}
      >
        Return to App</button-component
      >
    </div>`}};Ge.styles=[s`
      #header {
        --red: #c22a19;
      }
      #header h1 {
        font-size: 1.5rem;
        font-weight: 600;
        display: flex;
        gap: 1rem;
        color: var(--red);
      }

      #header p#intro {
        font-weight: 300;
        margin-top: 0.5rem;
        font-size: 0.85rem;
        color: #464646;
      }

      a#back {
        gap: 0.5rem;
        color: black;
        display: flex;
        align-items: center;
        text-decoration: none;
        font-weight: 300;
      }

      ui-icon {
        --primary-600: var(--red);
      }

      hr {
        border: 0.5px solid #dcdcdc;
        margin: 2rem 0;
      }
    `],Je([de()],Ge.prototype,"errorTitle",void 0),Je([de()],Ge.prototype,"errorDescription",void 0),Ge=Je([ae("locksmith-error")],Ge);var Ze=function(e,t,i,o){var r,s=arguments.length,a=s<3?t:null===o?o=Object.getOwnPropertyDescriptor(t,i):o;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(e,t,i,o);else for(var n=e.length-1;n>=0;n--)(r=e[n])&&(a=(s<3?r(a):s>3?r(t,i,a):r(t,i))||a);return s>3&&a&&Object.defineProperty(t,i,a),a};let Qe=class extends re{constructor(){super(...arguments),this.originOverride="",this.forceEmail="",this.inviteCode="",this.appName="FILL ME",this.minimumPasswordLength=6,this.showingPassword=!1,this.errorMsg=void 0,this.confirmEmailRequired=!1,this.validationok=!1,this.signUpRef=Pe(),this.emailRef=Pe(),this.passwordRef=Pe(),this.passwordConfirmationRef=Pe()}firstUpdated(){var e;null===(e=this.emailRef.value)||void 0===e||e.focus()}canSignIn(){var e,t,i;return(null===(e=this.emailRef.value)||void 0===e?void 0:e.value.length)>0&&(null===(t=this.passwordRef.value)||void 0===t?void 0:t.value.length)>0&&(null===(i=this.passwordConfirmationRef.value)||void 0===i?void 0:i.value.length)>0}doPasswordsMatch(){var e,t;return(null===(e=this.passwordRef.value)||void 0===e?void 0:e.value)===(null===(t=this.passwordConfirmationRef.value)||void 0===t?void 0:t.value)}async sendRegistrationRequest(){var e;const t={username:this.emailRef.value.value,email:this.emailRef.value.value,password:this.passwordRef.value.value,code:this.inviteCode,validationok:this.validationok};this.validationok=!1,this.confirmEmailRequired=!1,this.didYouMean=void 0;const i=await fetch(`${null!==(e=this.originOverride)&&void 0!==e?e:""}/api/register`,{method:"POST",body:JSON.stringify(t)});if(200!==i.status){if(409===i.status)throw new Error("This email is already being used.");if(400===i.status){const e=await i.json();if("password too short"===e.error)throw new Error("Password too short.");if("illegal username characters"===e.error)throw new Error("Email must be a valid email.");if(e.rejectEmail)throw new Error("This email address is invalid. If you need help, please contact support.");if(e.confirmEmail)throw e.didYouMean&&(this.didYouMean=e.didYouMean),this.confirmEmailRequired=!0,this.validationok=!0,new Error("We couldn't verify this email address. Please double-check for typos before trying again.")}throw new Error("Something went wrong.")}}passwordLongEnough(){var e;return(null===(e=this.passwordRef.value)||void 0===e?void 0:e.value.length)>=this.minimumPasswordLength}async attemptRegistration(){var e,t,i,o,r,s,a;if(this.errorMsg=void 0,!this.canSignIn())return this.errorMsg="Please enter a username and password.",void(0===(null===(e=this.emailRef.value)||void 0===e?void 0:e.value.length)?null===(t=this.emailRef.value)||void 0===t||t.focus():0===(null===(i=this.passwordRef.value)||void 0===i?void 0:i.value.length)?null===(o=this.passwordRef.value)||void 0===o||o.focus():null===(r=this.passwordConfirmationRef.value)||void 0===r||r.focus());if(!this.doPasswordsMatch())return this.errorMsg="The password must match.",void(null===(s=this.passwordConfirmationRef.value)||void 0===s||s.focus());if(!this.passwordLongEnough())return this.errorMsg=`Password must be at least ${this.minimumPasswordLength} characters long.`,void(null===(a=this.passwordRef.value)||void 0===a||a.focus());this.signUpRef.value.loading=!0,this.requestUpdate();try{await this.sendRegistrationRequest(),window.location.href="/login?onboard=true"}catch(e){console.error(e),this.errorMsg=e.message}finally{this.signUpRef.value.loading=!1}}render(){return D` <div id="root">
      <div id="header">
        <h1>Sign up to ${this.appName}</h1>
        ${0===this.forceEmail.length?D`
              <p id="intro">
                Already have an account? <a href="/login">Sign in instead</a>
              </p>
            `:D``}
        <p id="error">${this.errorMsg}</p>
      </div>

      <div id="inputs">
        <div class="input-container">
          <label for="username">Email Address</label>
          ${this.didYouMean?D`<button
                id="didYouMean"
                @click=${()=>{this.emailRef.value.value=this.didYouMean,this.emailRef.value.focus(),this.didYouMean=void 0}}
              >
                Did you mean <span>${this.didYouMean}</span>?
              </button>`:D``}
          <input
            id="username"
            ${Me(this.emailRef)}
            autofill="username"
            autocapitalize="off"
            autocapitalize="off"
            placeholder="Your email"
            value="${this.forceEmail}"
            ?disabled=${this.forceEmail.length>0}
            @input=${()=>{this.validationok=!1,this.didYouMean=void 0}}
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
            ${Me(this.passwordRef)}
            autocomplete="new-password"
            type="${this.showingPassword?"text":"password"}"
            placeholder="Your password"
          />
        </div>

        <div class="input-container">
          <label for="password">Confirm your Password</label>
          <input
            id="password"
            ${Me(this.passwordConfirmationRef)}
            autocomplete="new-password"
            type="${this.showingPassword?"text":"password"}"
            placeholder="Confirm your Password"
          />
        </div>
      </div>

      <button-component
        ${Me(this.signUpRef)}
        class="big"
        .expectLoad=${!0}
        .loadingText=${"Signing Up.."}
        @fl-click=${this.attemptRegistration}
        >${this.confirmEmailRequired?"Confirm & Sign Up":"Sign Up"}</button-component
      >
    </div>`}};Qe.styles=[he,s`
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

      button#didYouMean {
        background: none;
        border: none;
        color: #b8123a;
        font-size: 0.85rem;
        text-align: left;
        margin-bottom: 0.5rem;
        padding: 0;
      }

      button#didYouMean > span {
        text-decoration: underline;
      }

      ui-icon {
        --primary-600: var(--accent);
        --primary-500: var(--accent);
      }
    `],Ze([de()],Qe.prototype,"originOverride",void 0),Ze([de()],Qe.prototype,"forceEmail",void 0),Ze([de()],Qe.prototype,"inviteCode",void 0),Ze([de()],Qe.prototype,"appName",void 0),Ze([de()],Qe.prototype,"minimumPasswordLength",void 0),Ze([ce()],Qe.prototype,"showingPassword",void 0),Ze([ce()],Qe.prototype,"errorMsg",void 0),Ze([ce()],Qe.prototype,"didYouMean",void 0),Ze([ce()],Qe.prototype,"confirmEmailRequired",void 0),Ze([ce()],Qe.prototype,"validationok",void 0),Qe=Ze([ae("locksmith-registration")],Qe);var Xe=function(e,t,i,o){var r,s=arguments.length,a=s<3?t:null===o?o=Object.getOwnPropertyDescriptor(t,i):o;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(e,t,i,o);else for(var n=e.length-1;n>=0;n--)(r=e[n])&&(a=(s<3?r(a):s>3?r(t,i,a):r(t,i))||a);return s>3&&a&&Object.defineProperty(t,i,a),a};let et=class extends re{constructor(){super(...arguments),this.hasResetCode="",this.appName="FILL ME",this.originOverride="",this.minimumPasswordLength=6,this.showingPassword=!1,this.errorMsg=void 0,this.loadedOnce=!1,this.stage=0,this.passwordRef=Pe(),this.passwordConfirmationRef=Pe(),this.resetButtonRef=Pe(),this.emailRef=Pe(),this.resetFullyButtonRef=Pe()}firstUpdated(){var e,t;null===(e=this.emailRef.value)||void 0===e||e.focus(),null===(t=this.passwordRef.value)||void 0===t||t.focus(),setTimeout((()=>{this.loadedOnce=!0}))}async attemptReset(){var e,t;if(0===(null===(e=this.emailRef.value)||void 0===e?void 0:e.value.length))return void(this.errorMsg="Please enter an email address.");this.resetButtonRef.value.loading=!0;if(200!==(await fetch(`${null!==(t=this.originOverride)&&void 0!==t?t:""}/api/reset-password?username=${this.emailRef.value.value}`,{method:"POST"})).status)return console.error("Something bad happened while resetting a password"),void(this.errorMsg="Something went wrong. Please try again later.");this.stage=1}doPasswordsMatch(){var e,t;return(null===(e=this.passwordRef.value)||void 0===e?void 0:e.value)===(null===(t=this.passwordConfirmationRef.value)||void 0===t?void 0:t.value)}passwordLongEnough(){var e;return(null===(e=this.passwordRef.value)||void 0===e?void 0:e.value.length)>=this.minimumPasswordLength}async fullyResetPassword(){var e,t;if(!this.doPasswordsMatch())return void(this.errorMsg="The password must match.");if(!this.passwordLongEnough())return this.errorMsg=`Password must be at least ${this.minimumPasswordLength} characters long.`,void(null===(e=this.passwordRef.value)||void 0===e||e.focus());this.resetFullyButtonRef.value.loading=!0;if(200!==(await fetch(`${null!==(t=this.originOverride)&&void 0!==t?t:""}/api/reset-password`,{method:"PATCH",body:JSON.stringify({password:this.passwordRef.value.value})})).status)return console.error("Something bad happened while resetting a password"),void(this.errorMsg="Something went wrong. Please try again later.");this.stage=2}render(){return 2==this.stage?D`<div id="root" class="${this.loadedOnce?"":"hide"}">
        <div id="header">
          <h1>Your password has been reset.</h1>
          <p id="intro">
            <strong>You may now login using the new password.</strong> If you
            need any help, feel free to contact us.
          </p>
        </div>

        <a href="/login">Go to Login</a>
      </div>`:""!==this.hasResetCode?D`<div id="root" class="${this.loadedOnce?"":"hide"}">
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
              ${Me(this.passwordRef)}
              autocomplete="new-password"
              type="${this.showingPassword?"text":"password"}"
              placeholder="Your password"
            />
          </div>

          <div class="input-container">
            <label for="password">Confirm your Password</label>
            <input
              id="password"
              ${Me(this.passwordConfirmationRef)}
              autocomplete="new-password"
              type="${this.showingPassword?"text":"password"}"
              placeholder="Confirm your Password"
            />
          </div>
        </div>
        <button-component
          ${Me(this.resetFullyButtonRef)}
          class="big"
          .expectLoad=${!0}
          .loadingText=${"Resetting.."}
          @fl-click=${this.fullyResetPassword}
          >Reset Password</button-component
        >
      </div>`:1==this.stage?D` <div id="root" class="${this.loadedOnce?"":"hide"}">
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
      </div>`:D` <div id="root" class="${this.loadedOnce?"":"hide"}">
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
            ${Me(this.emailRef)}
            autofill="username"
            autocapitalize="off"
            autocapitalize="off"
            placeholder="Your email"
          />
        </div>
      </div>

      <button-component
        ${Me(this.resetButtonRef)}
        class="big"
        .expectLoad=${!0}
        .loadingText=${"Sending.."}
        @fl-click=${this.attemptReset}
        >Send Reset Link</button-component
      >
    </div>`}};et.styles=[he,s`
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
    `],Xe([de()],et.prototype,"hasResetCode",void 0),Xe([de()],et.prototype,"appName",void 0),Xe([de()],et.prototype,"originOverride",void 0),Xe([de()],et.prototype,"minimumPasswordLength",void 0),Xe([ce()],et.prototype,"showingPassword",void 0),Xe([ce()],et.prototype,"errorMsg",void 0),Xe([ce()],et.prototype,"loadedOnce",void 0),Xe([ce()],et.prototype,"stage",void 0),et=Xe([ae("locksmith-reset-password")],et);var tt=function(e,t,i,o){var r,s=arguments.length,a=s<3?t:null===o?o=Object.getOwnPropertyDescriptor(t,i):o;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(e,t,i,o);else for(var n=e.length-1;n>=0;n--)(r=e[n])&&(a=(s<3?r(a):s>3?r(t,i,a):r(t,i))||a);return s>3&&a&&Object.defineProperty(t,i,a),a};let it=class extends re{constructor(){super(...arguments),this.pages=[{URLKey:"account",Name:"Account",TemplateLiteral:D`<p>Hello, Account page</p>`},{URLKey:"security",Name:"Security",TemplateLiteral:D`<p>Hello, Security page</p>`},{URLKey:"demo",Name:"Demo",TemplateLiteral:D`<p>Hello, Demo page</p>`}],this.defaultPageKey="demo",this.selected="",this.activePageComponent=void 0}async setSelected(e,t=!0){var i;if(this.selected=e.URLKey,t&&history.pushState({key:e.URLKey},"",`#${e.URLKey}`),void 0!==e.PageComponent){const t=new e.PageComponent;if(void 0!==e.LoadProps){const i=e.LoadProps();t.setProps(i)}return t.OnPageLoad&&await t.OnPageLoad(),void(this.activePageComponent=t)}this.activePageComponent=null!==(i=e.TemplateLiteral)&&void 0!==i?i:D`<p>Page is missing it's definition.</p.>`}firstUpdated(){var e;const t=(null===(e=location.hash)||void 0===e?void 0:e.replace("#",""))||this.defaultPageKey;let i=this.pages.find((e=>e.URLKey===t));i||(i=this.pages.find((e=>e.URLKey===this.defaultPageKey))),i?this.setSelected(i,!1).then((()=>{this.updateComplete.then((()=>{const e=this.renderRoot.querySelector("button.selected");null==e||e.scrollIntoView({behavior:"smooth",inline:"center",block:"nearest"})}))})):console.warn("Page key not found: ",this.defaultPageKey)}connectedCallback(){super.connectedCallback(),window.addEventListener("popstate",this.listenForNavChanges.bind(this))}disconnectedCallback(){super.connectedCallback(),window.removeEventListener("popstate",this.listenForNavChanges.bind(this))}listenForNavChanges(){const e=location.hash.replace("#","")||this.defaultPageKey,t=this.pages.find((t=>t.URLKey===e));t&&(this.setSelected(t,!1),this.updateComplete.then((()=>{const e=this.renderRoot.querySelector("button.selected");null==e||e.scrollIntoView({behavior:"smooth",inline:"center",block:"nearest"})})))}render(){return D`<nav>
        ${this.pages.map((e=>D`<button
              class=${ge({selected:this.selected===e.URLKey})}
              @click=${()=>this.setSelected(e)}
            >
              ${e.Name}
            </button>`))}
      </nav>
      <section>
        ${void 0!==this.activePageComponent?this.activePageComponent:""}
      </section>`}};it.styles=s`
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
  `,tt([de()],it.prototype,"pages",void 0),tt([de()],it.prototype,"defaultPageKey",void 0),tt([ce()],it.prototype,"selected",void 0),tt([ce()],it.prototype,"activePageComponent",void 0),it=tt([ae("quick-nav")],it);var ot=function(e,t,i,o){var r,s=arguments.length,a=s<3?t:null===o?o=Object.getOwnPropertyDescriptor(t,i):o;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(e,t,i,o);else for(var n=e.length-1;n>=0;n--)(r=e[n])&&(a=(s<3?r(a):s>3?r(t,i,a):r(t,i))||a);return s>3&&a&&Object.defineProperty(t,i,a),a};let rt=class extends re{constructor(){super(...arguments),this.aboutMeState=new Be(this,Fe)}render(){return D`<div class="input-container">
      <label>
        Email Address
        ${Fe.hasPermission("user.update.email")?D` <button>Change Email</button> `:D``}
      </label>
      ${Fe.hasPermission("user.update.email")?D``:D` <p>Please contact us to change your account email.</p> `}
      <input value="${Fe.email}" disabled />
    </div>`}};rt.styles=[he,s`
      * {
        box-sizing: border-box;
      }
    `],rt=ot([ae("locksmith-update-email")],rt);var st=function(e,t,i,o){var r,s=arguments.length,a=s<3?t:null===o?o=Object.getOwnPropertyDescriptor(t,i):o;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(e,t,i,o);else for(var n=e.length-1;n>=0;n--)(r=e[n])&&(a=(s<3?r(a):s>3?r(t,i,a):r(t,i))||a);return s>3&&a&&Object.defineProperty(t,i,a),a};let at=class extends re{constructor(){super(...arguments),this.aboutMeState=new Be(this,Fe)}connectedCallback(){super.connectedCallback(),Fe.loadIfNeeded(!0)}render(){return D`<div>
      <a href="/app" id="back"
        ><ui-icon name="home" size="1.15rem"></ui-icon> Back to App</a
      >

      <div id="header">
        <h1>Hi, ${Fe.username}.</h1>
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
        .pages=${[{URLKey:"account",Name:"Account",PageComponent:rt},{URLKey:"security",Name:"Security",TemplateLiteral:D`<p>TODO</p>`}]}
      ></quick-nav>
    </div>`}};at.styles=[he,s`
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
    `],at=st([ae("locksmith-profile")],at);var nt=function(e,t,i,o){var r,s=arguments.length,a=s<3?t:null===o?o=Object.getOwnPropertyDescriptor(t,i):o;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(e,t,i,o);else for(var n=e.length-1;n>=0;n--)(r=e[n])&&(a=(s<3?r(a):s>3?r(t,i,a):r(t,i))||a);return s>3&&a&&Object.defineProperty(t,i,a),a};let lt=class extends re{constructor(){super(),this.resendButtonRef=Pe(),this.originOverride="",this.email="",this.errorMsg=void 0,this.didResend=!1,this.loadedOnce=!1,this.resendDisabled=!0,this.resendSecondsLeft=0,this.startResendPeriod(30)}startResendPeriod(e){this.resendDisabled=!0,this.resendSecondsLeft=e;const t=setInterval((()=>{this.resendSecondsLeft--,this.resendSecondsLeft<=0&&(clearInterval(t),this.resendDisabled=!1)}),1e3)}async resendLink(){var e;if(!this.resendDisabled)try{this.resendButtonRef.value.loading=!0;const t=await fetch(`${null!==(e=this.originOverride)&&void 0!==e?e:""}/api/verify/resend`,{method:"POST"});if(this.startResendPeriod(60),200!==t.status)return 429===t.status?void(this.errorMsg="Please wait before trying again or contact support."):void(this.errorMsg="Please contact support.");this.didResend=!0}catch(e){console.error(e),this.errorMsg="Please contact support."}finally{this.resendButtonRef.value.loading=!1}}render(){return D`<div id="root">
      <div id="header">
        <h1>Please verify your email.</h1>
        ${this.errorMsg?D`<p
              id="error"
              aria-live="assertive"
              role="status"
              aria-atomic="true"
              aria-relevant="additions"
            >
              ${this.errorMsg}
            </p>`:D``}
        ${this.didResend?D`<p
              id="resend"
              aria-live="assertive"
              role="status"
              aria-atomic="true"
              aria-relevant="additions"
            >
              Verification resent successfully.
            </p>`:D``}
        <p id="intro">
          We've sent an email to the following email address. Please click the
          verification link in the email to continue.<br /><br />

          <strong>Didn't get the email?</strong> Make sure to check your spam
          folder.<br /><br />

          <strong>Mistype the email?</strong> Please contact support and we'll
          get it changed.
        </p>
      </div>

      <p id="email">${this.email}</p>

      <button-component
        ${Me(this.resendButtonRef)}
        class="big"
        .expectLoad=${!0}
        .disabled=${this.resendDisabled}
        .loadingText=${"Sending.."}
        @fl-click=${this.resendLink}
      >
        ${this.resendDisabled?`Resend in ${this.resendSecondsLeft}s...`:"Resend Link"}</button-component
      >
    </div>`}};lt.styles=[he,s`
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

      p#resend {
        color: var(--accent);
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

      #email {
        font-size: 1rem;
        padding: 0.75rem;
        background-color: var(--gray-50, #f8f8f8);
        border: 1px solid var(--gray-200, #dcdcdc);
        border-radius: 0.25rem;
      }
    `],nt([de()],lt.prototype,"originOverride",void 0),nt([de()],lt.prototype,"email",void 0),nt([ce()],lt.prototype,"errorMsg",void 0),nt([ce()],lt.prototype,"didResend",void 0),nt([ce()],lt.prototype,"loadedOnce",void 0),nt([ce()],lt.prototype,"resendDisabled",void 0),nt([ce()],lt.prototype,"resendSecondsLeft",void 0),lt=nt([ae("locksmith-verify-email")],lt);var dt=function(e,t,i,o){var r,s=arguments.length,a=s<3?t:null===o?o=Object.getOwnPropertyDescriptor(t,i):o;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(e,t,i,o);else for(var n=e.length-1;n>=0;n--)(r=e[n])&&(a=(s<3?r(a):s>3?r(t,i,a):r(t,i))||a);return s>3&&a&&Object.defineProperty(t,i,a),a};let ct=class extends re{constructor(){super(),this.originOverride="",this.code="",this.errorMsg=void 0,this.didVerify=!1,this.isExchanging=!0}firstUpdated(){this.exchangeCode()}async exchangeCode(){var e;this.isExchanging=!0,this.errorMsg=void 0,this.didVerify=!1;try{const t=await fetch(`${null!==(e=this.originOverride)&&void 0!==e?e:""}/api/verify/exchange`,{method:"POST",headers:{"content-type":"application/json"},body:JSON.stringify({code:this.code})});if(!t.ok||200!==t.status)return void(429===t.status?this.errorMsg="This verification link has expired. Please request a new one.":this.errorMsg="Failed to verify this link. Please try verification again.");window.location.href="/app"}catch(e){console.error(e),this.errorMsg="Please contact support."}finally{this.isExchanging=!1}}getTitle(){return this.errorMsg?"Failed to verify":this.didVerify?"Verification successful":"Verifying"}render(){return this.isExchanging?D`
        <div
          id="loadingState"
          aria-live="polite"
          role="status"
          aria-atomic="true"
        >
          <div id="loadingCard">
            <div class="spinner" aria-hidden="true"></div>
            <p id="loadingTitle">Verifying your account</p>
            <p id="loadingText">
              Please wait a moment while we verify your account.
            </p>
          </div>
        </div>
      `:D`
      <div id="root">
        <div id="header">
          <h1>${this.getTitle()}.</h1>

          ${this.errorMsg?D`
                <p
                  id="error"
                  aria-live="assertive"
                  role="status"
                  aria-atomic="true"
                  aria-relevant="additions"
                >
                  ${this.errorMsg}
                </p>
              `:D``}
          ${this.didVerify?D`
                <p
                  id="success"
                  aria-live="polite"
                  role="status"
                  aria-atomic="true"
                >
                  You'll be redirected momentarily...
                </p>
              `:D``}

          <p id="intro">
            ${this.didVerify?"Your email has been verified successfully.":"We couldn't verify this link."}
          </p>
        </div>
      </div>
    `}};ct.styles=[he,s`
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
        flex-direction: column;
        gap: 1rem;
        opacity: 1;
      }

      #header h1 {
        font-size: 1.5rem;
        font-weight: 500;
      }

      #header p#intro {
        font-weight: 300;
        margin-top: 0.5rem;
        font-size: 0.9rem;
        color: #464646;
        line-height: 1.35rem;
      }

      p#error {
        color: #b8123a;
        margin-top: 0.5rem;
      }

      p#success {
        color: var(--accent);
        margin-top: 0.5rem;
      }

      #loadingState {
        min-height: 220px;
        display: flex;
        align-items: center;
        justify-content: center;
      }

      #loadingCard {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        gap: 1rem;
        text-align: center;
      }

      .spinner {
        width: 2.5rem;
        height: 2.5rem;
        border-radius: 999px;
        border: 3px solid rgba(0, 0, 0, 0.12);
        border-top-color: var(--accent);
        animation: spin 0.8s linear infinite;
        flex-shrink: 0;
      }

      #loadingTitle {
        font-size: 1.1rem;
        font-weight: 500;
        color: #1f1f1f;
      }

      #loadingText {
        font-size: 0.92rem;
        line-height: 1.35rem;
        color: #5a5a5a;
        max-width: 24rem;
      }

      @keyframes spin {
        to {
          transform: rotate(360deg);
        }
      }

      ui-icon {
        --primary-600: var(--accent);
        --primary-500: var(--accent);
      }
    `],dt([de()],ct.prototype,"originOverride",void 0),dt([de()],ct.prototype,"code",void 0),dt([ce()],ct.prototype,"errorMsg",void 0),dt([ce()],ct.prototype,"didVerify",void 0),dt([ce()],ct.prototype,"isExchanging",void 0),ct=dt([ae("locksmith-verify-code")],ct);var ht=function(e,t,i,o){var r,s=arguments.length,a=s<3?t:null===o?o=Object.getOwnPropertyDescriptor(t,i):o;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(e,t,i,o);else for(var n=e.length-1;n>=0;n--)(r=e[n])&&(a=(s<3?r(a):s>3?r(t,i,a):r(t,i))||a);return s>3&&a&&Object.defineProperty(t,i,a),a};let pt=class extends re{render(){return D`<div id="root">
      <header>
        ${void 0!==this.logoURL&&""!==this.logoURL?D` <img src="${this.logoURL}" /> `:D``}
      </header>
      <main>
        <div id="slotWrapper">
          <slot name="main"></slot>
        </div>
      </main>
      <footer></footer>
    </div>`}};async function ut(e,t){return fetch(e,t)}pt.styles=[s`
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
    `],ht([de()],pt.prototype,"logoURL",void 0),pt=ht([ae("locksmith-layout")],pt);var vt=function(e,t,i,o){var r,s=arguments.length,a=s<3?t:null===o?o=Object.getOwnPropertyDescriptor(t,i):o;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(e,t,i,o);else for(var n=e.length-1;n>=0;n--)(r=e[n])&&(a=(s<3?r(a):s>3?r(t,i,a):r(t,i))||a);return s>3&&a&&Object.defineProperty(t,i,a),a};let ft=class extends re{constructor(){super(...arguments),this.location={vertical:"top",horizontal:"right"},this.open=!1,this.launchpadForceClosed=!1,this.aboutMeState=new Be(this,Fe),this.manageAccountButton=Pe(),this.handleEscape=e=>{"Escape"===e.key&&(this.open=!1)},this.listenForOOBClicks=e=>{if(!this.shadowRoot)return;e.composedPath().includes(this.shadowRoot.host)||(this.open=!1,window.removeEventListener("click",this.listenForOOBClicks))}}updated(){this.setAttribute("location-vertical",this.location.vertical),this.setAttribute("location-horizontal",this.location.horizontal)}connectedCallback(){super.connectedCallback(),Fe.loadIfNeeded()}openClicked(){this.open=!this.open,setTimeout((()=>{this.open?(window.addEventListener("click",this.listenForOOBClicks),window.addEventListener("keydown",this.handleEscape),this.updateComplete.then((()=>{setTimeout((()=>{var e;null===(e=this.manageAccountButton.value)||void 0===e||e.focus()}),100)}))):(window.removeEventListener("click",this.listenForOOBClicks),window.removeEventListener("keydown",this.handleEscape))}))}render(){return D` <div id="container">
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
              <p id="title">${Fe.username}</p>
              <p id="desc">${Fe.email}</p>
            </div>
          </div>

          <div id="actions">
            <button
              ${Me(this.manageAccountButton)}
              @click=${()=>{window.location.href="/profile"}}
            >
              <ui-icon name="cog" size="1rem"></ui-icon>
              Manage Account
              <span></span>
            </button>
            ${Fe.hasPermission("view.ls-admin")?D`
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
              @click=${()=>{Fe.signOut()}}
            >
              <ui-icon name="sign-out" size="1rem" colorway="danger"></ui-icon>
              Sign out
            </button>
          </div>
        </div>
      </div>

      ${void 0===Fe.isLaunchpad||this.launchpadForceClosed?void 0:D`<div id="launchpad">
            <div>
              <p id="launchpad-status">
                ${Fe.isLaunchpad.toUpperCase()} Role
              </p>
              <p id="launchpad-user">
                This is what the ${Fe.isLaunchpad.toUpperCase()} role would
                see.
              </p>
            </div>

            <div>
              <button
                @click=${()=>{this.launchpadForceClosed=!0}}
              >
                &times;
              </button>
            </div>
          </div>`}`}};ft.styles=[s`
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
        background-color: var(--danger-600, #e01e47);
        color: var(--danger-50, #fff1f2);
        font-size: 1rem;
        z-index: 1000;

        border-top: 2px solid var(--danger-800, #9f1239);

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
        color: var(--danger-50, #fff1f2);
      }

      #launchpad button {
        color: var(--danger-50, #fff1f2);
        margin: 0;
        padding: 0;
        border: none;
        background-color: transparent;
        font-size: 1.5rem;
        cursor: pointer;
      }
    `],vt([de()],ft.prototype,"location",void 0),vt([ce()],ft.prototype,"open",void 0),vt([ce()],ft.prototype,"launchpadForceClosed",void 0),ft=vt([ae("locksmith-user-icon")],ft);export{qe as AboutMeState,Ee as GenerateFingerprint,Ge as LocksmithErrorComponent,pt as LocksmithLayout,Ke as LocksmithLoginComponent,at as LocksmithProfileComponent,Qe as LocksmithRegistrationComponent,et as LocksmithResetPasswordComponent,ft as LocksmithUserIconComponent,ct as LocksmithVerifyCodeComponent,lt as LocksmithVerifyEmailComponent,ut as SecureFetch,Fe as aboutMe,he as inputStyles};
//# sourceMappingURL=locksmith-ui.bundle.js.map
