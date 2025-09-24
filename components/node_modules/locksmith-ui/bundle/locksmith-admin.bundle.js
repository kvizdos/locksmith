function e(e,t){void 0===t&&(t={});for(var r=function(e){for(var t=[],r=0;r<e.length;){var i=e[r];if("*"!==i&&"+"!==i&&"?"!==i)if("\\"!==i)if("{"!==i)if("}"!==i)if(":"!==i)if("("!==i)t.push({type:"CHAR",index:r,value:e[r++]});else{var o=1,n="";if("?"===e[s=r+1])throw new TypeError('Pattern cannot start with "?" at '.concat(s));for(;s<e.length;)if("\\"!==e[s]){if(")"===e[s]){if(0==--o){s++;break}}else if("("===e[s]&&(o++,"?"!==e[s+1]))throw new TypeError("Capturing groups are not allowed at ".concat(s));n+=e[s++]}else n+=e[s++]+e[s++];if(o)throw new TypeError("Unbalanced pattern at ".concat(r));if(!n)throw new TypeError("Missing pattern at ".concat(r));t.push({type:"PATTERN",index:r,value:n}),r=s}else{for(var a="",s=r+1;s<e.length;){var l=e.charCodeAt(s);if(!(l>=48&&l<=57||l>=65&&l<=90||l>=97&&l<=122||95===l))break;a+=e[s++]}if(!a)throw new TypeError("Missing parameter name at ".concat(r));t.push({type:"NAME",index:r,value:a}),r=s}else t.push({type:"CLOSE",index:r,value:e[r++]});else t.push({type:"OPEN",index:r,value:e[r++]});else t.push({type:"ESCAPED_CHAR",index:r++,value:e[r++]});else t.push({type:"MODIFIER",index:r,value:e[r++]})}return t.push({type:"END",index:r,value:""}),t}(e),o=t.prefixes,n=void 0===o?"./":o,a=t.delimiter,s=void 0===a?"/#?":a,l=[],c=0,d=0,p="",h=function(e){if(d<r.length&&r[d].type===e)return r[d++].value},u=function(e){var t=h(e);if(void 0!==t)return t;var i=r[d],o=i.type,n=i.index;throw new TypeError("Unexpected ".concat(o," at ").concat(n,", expected ").concat(e))},m=function(){for(var e,t="";e=h("CHAR")||h("ESCAPED_CHAR");)t+=e;return t},f=function(e){var t=l[l.length-1],r=e||(t&&"string"==typeof t?t:"");if(t&&!r)throw new TypeError('Must have text between two parameters, missing text after "'.concat(t.name,'"'));return!r||function(e){for(var t=0,r=s;t<r.length;t++){var i=r[t];if(e.indexOf(i)>-1)return!0}return!1}(r)?"[^".concat(i(s),"]+?"):"(?:(?!".concat(i(r),")[^").concat(i(s),"])+?")};d<r.length;){var v=h("CHAR"),g=h("NAME"),y=h("PATTERN");if(g||y){var b=v||"";-1===n.indexOf(b)&&(p+=b,b=""),p&&(l.push(p),p=""),l.push({name:g||c++,prefix:b,suffix:"",pattern:y||f(b),modifier:h("MODIFIER")||""})}else{var w=v||h("ESCAPED_CHAR");if(w)p+=w;else if(p&&(l.push(p),p=""),h("OPEN")){b=m();var x=h("NAME")||"",$=h("PATTERN")||"",k=m();u("CLOSE"),l.push({name:x||($?c++:""),pattern:x&&!$?f(b):$,prefix:b,suffix:k,modifier:h("MODIFIER")||""})}else u("END")}}return l}function t(t,i){return r(e(t,i),i)}function r(e,t){void 0===t&&(t={});var r=o(t),i=t.encode,n=void 0===i?function(e){return e}:i,a=t.validate,s=void 0===a||a,l=e.map((function(e){if("object"==typeof e)return new RegExp("^(?:".concat(e.pattern,")$"),r)}));return function(t){for(var r="",i=0;i<e.length;i++){var o=e[i];if("string"!=typeof o){var a=t?t[o.name]:void 0,c="?"===o.modifier||"*"===o.modifier,d="*"===o.modifier||"+"===o.modifier;if(Array.isArray(a)){if(!d)throw new TypeError('Expected "'.concat(o.name,'" to not repeat, but got an array'));if(0===a.length){if(c)continue;throw new TypeError('Expected "'.concat(o.name,'" to not be empty'))}for(var p=0;p<a.length;p++){var h=n(a[p],o);if(s&&!l[i].test(h))throw new TypeError('Expected all "'.concat(o.name,'" to match "').concat(o.pattern,'", but got "').concat(h,'"'));r+=o.prefix+h+o.suffix}}else if("string"!=typeof a&&"number"!=typeof a){if(!c){var u=d?"an array":"a string";throw new TypeError('Expected "'.concat(o.name,'" to be ').concat(u))}}else{h=n(String(a),o);if(s&&!l[i].test(h))throw new TypeError('Expected "'.concat(o.name,'" to match "').concat(o.pattern,'", but got "').concat(h,'"'));r+=o.prefix+h+o.suffix}}else r+=o}return r}}function i(e){return e.replace(/([.+*?=^!:${}()[\]|/\\])/g,"\\$1")}function o(e){return e&&e.sensitive?"":"i"}function n(t,r,n){return function(e,t,r){void 0===r&&(r={});for(var n=r.strict,a=void 0!==n&&n,s=r.start,l=void 0===s||s,c=r.end,d=void 0===c||c,p=r.encode,h=void 0===p?function(e){return e}:p,u=r.delimiter,m=void 0===u?"/#?":u,f=r.endsWith,v="[".concat(i(void 0===f?"":f),"]|$"),g="[".concat(i(m),"]"),y=l?"^":"",b=0,w=e;b<w.length;b++){var x=w[b];if("string"==typeof x)y+=i(h(x));else{var $=i(h(x.prefix)),k=i(h(x.suffix));if(x.pattern)if(t&&t.push(x),$||k)if("+"===x.modifier||"*"===x.modifier){var _="*"===x.modifier?"?":"";y+="(?:".concat($,"((?:").concat(x.pattern,")(?:").concat(k).concat($,"(?:").concat(x.pattern,"))*)").concat(k,")").concat(_)}else y+="(?:".concat($,"(").concat(x.pattern,")").concat(k,")").concat(x.modifier);else{if("+"===x.modifier||"*"===x.modifier)throw new TypeError('Can not repeat "'.concat(x.name,'" without a prefix and suffix'));y+="(".concat(x.pattern,")").concat(x.modifier)}else y+="(?:".concat($).concat(k,")").concat(x.modifier)}}if(d)a||(y+="".concat(g,"?")),y+=r.endsWith?"(?=".concat(v,")"):"$";else{var A=e[e.length-1],C="string"==typeof A?g.indexOf(A[A.length-1])>-1:void 0===A;a||(y+="(?:".concat(g,"(?=").concat(v,"))?")),C||(y+="(?=".concat(g,"|").concat(v,")"))}return new RegExp(y,o(r))}(e(t,n),r,n)}function a(e,t,r){return e instanceof RegExp?function(e,t){if(!t)return e;for(var r=/\((?:\?<(.*?)>)?(?!\?)/g,i=0,o=r.exec(e.source);o;)t.push({name:o[1]||i++,prefix:"",suffix:"",modifier:"",pattern:""}),o=r.exec(e.source);return e}(e,t):Array.isArray(e)?function(e,t,r){var i=e.map((function(e){return a(e,t,r).source}));return new RegExp("(?:".concat(i.join("|"),")"),o(r))}(e,t,r):n(e,t,r)}function s(e){return"object"==typeof e&&!!e}function l(e){return"function"==typeof e}function c(e){return"string"==typeof e}function d(e=[]){return Array.isArray(e)?e:[e]}function p(e){return`[Vaadin.Router] ${e}`}class h extends Error{code;context;constructor(e){super(p(`Page not found (${e.pathname})`)),this.context=e,this.code=404}}const u=Symbol("NotFoundResult");function m(e){return new h(e)}function f(e){return(Array.isArray(e)?e[0]:e)??""}function v(e){return f(e?.path)}const g=new Map;function y(e){try{return decodeURIComponent(e)}catch{return e}}g.set("|false",{keys:[],pattern:/(?:)/u});var b=function(e,t,r=!1,i=[],o){const n=`${e}|${String(r)}`,s=f(t);let l=g.get(n);if(!l){const t=[];l={keys:t,pattern:a(e,t,{end:r,strict:""===e})},g.set(n,l)}const c=l.pattern.exec(s);if(!c)return null;const d={...o};for(let e=1;e<c.length;e++){const t=l.keys[e-1],r=t.name,i=c[e];void 0===i&&Object.hasOwn(d,r)||("+"===t.modifier||"*"===t.modifier?d[r]=i?i.split(/[/?#]/u).map(y):[]:d[r]=i?y(i):i)}return{keys:[...i,...l.keys],params:d,path:c[0]}};var w=function e(t,r,i,o,n){let a,s,l=0,c=v(t);return c.startsWith("/")&&(i&&(c=c.substring(1)),i=!0),{next(d){if(t===d)return{done:!0,value:void 0};t.__children??=function(e){return Array.isArray(e)&&e.length>0?e:void 0}(t.children);const p=t.__children??[],h=!t.__children&&!t.children;if(!a&&(a=b(c,r,h,o,n),a))return{value:{keys:a.keys,params:a.params,path:a.path,route:t}};if(a&&p.length>0)for(;l<p.length;){if(!s){const o=p[l];o.parent=t;let n=a.path.length;n>0&&"/"===r.charAt(n)&&(n+=1),s=e(o,r.substring(n),i,a.keys,a.params)}const o=s.next(d);if(!o.done)return{done:!1,value:o.value};s=null,l+=1}return{done:!0,value:void 0}}}};function x(e){if(l(e.route.action))return e.route.action(e)}class $ extends Error{code;context;constructor(e,t){let r=`Path '${e.pathname}' is not properly resolved due to an error.`;const i=v(e.route);i&&(r+=` Resolution had failed on route: '${i}'`),super(r,t),this.code=t?.code,this.context=e}warn(){console.warn(this.message)}}class k{baseUrl;#e;errorHandler;resolveRoute;#t;constructor(e,{baseUrl:t="",context:r,errorHandler:i,resolveRoute:o=x}={}){if(Object(e)!==e)throw new TypeError("Invalid routes");this.baseUrl=t,this.errorHandler=i,this.resolveRoute=o,Array.isArray(e)?this.#t={__children:e,__synthetic:!0,action:()=>{},path:""}:this.#t={...e,parent:void 0},this.#e={...r,hash:"",next:async()=>u,params:{},pathname:"",resolver:this,route:this.#t,search:"",chain:[]}}get root(){return this.#t}get context(){return this.#e}get __effectiveBaseUrl(){return this.baseUrl?new URL(this.baseUrl,document.baseURI||document.URL).href.replace(/[^/]*$/u,""):""}getRoutes(){return[...this.#t.__children??[]]}removeRoutes(){this.#t.__children=[]}async resolve(e){const t=this,r={...this.#e,...c(e)?{pathname:e}:e,next:l},i=w(this.#t,this.__normalizePathname(r.pathname)??r.pathname,!!this.baseUrl),o=this.resolveRoute;let n=null,a=null,s=r;async function l(e=!1,c=n?.value?.route,d){const p=null===d?n?.value?.route:void 0;if(n=a??i.next(p),a=null,!e&&(n.done||!function(e,t){let r=e;for(;r;)if(r=r.parent,r===t)return!0;return!1}(n.value.route,c)))return a=n,u;if(n.done)throw m(r);s={...r,params:n.value.params,route:n.value.route,chain:s.chain?.slice()},function(e,t){const{path:r,route:i}=t;if(i&&!i.__synthetic){const t={path:r,route:i};if(i.parent&&e.chain)for(let t=e.chain.length-1;t>=0&&e.chain[t].route!==i.parent;t--)e.chain.pop();e.chain?.push(t)}}(s,n.value);const h=await o(s);return null!=h&&h!==u?(s.result=(f=h)&&"object"==typeof f&&"next"in f&&"params"in f&&"result"in f&&"route"in f?h.result:h,t.#e=s,s):await l(e,c,h);var f}try{return await l(!0,this.#t)}catch(e){const t=e instanceof h?e:new $(s,{code:500,cause:e});if(this.errorHandler)return s.result=this.errorHandler(t),s;throw e}}setRoutes(e){this.#t.__children=[...d(e)]}__normalizePathname(e){if(!this.baseUrl)return e;const t=this.__effectiveBaseUrl,r=e.startsWith("/")?new URL(t).origin+e:`./${e}`,i=new URL(r,t).href;return i.startsWith(t)?i.slice(t.length):void 0}addRoutes(e){return this.#t.__children=[...this.#t.__children??[],...d(e)],this.getRoutes()}}function _(e,t,r,i){const o=t.name??i?.(t);if(o&&(e.has(o)?e.get(o)?.push(t):e.set(o,[t])),Array.isArray(r))for(const o of r)o.parent=t,_(e,o,o.__children??o.children,i)}function A(e,t){const r=e.get(t);if(r){if(r.length>1)throw new Error(`Duplicate route with name "${t}". Try seting unique 'name' route properties.`);return r[0]}}var C=function(t,i={}){if(!(t instanceof k))throw new TypeError("An instance of Resolver is expected");const o=new Map,n=new Map;return(a,s)=>{let l=A(n,a);if(!l&&(n.clear(),_(n,t.root,t.root.__children,i.cacheKeyProvider),l=A(n,a),!l))throw new Error(`Route "${a}" not found`);let d=l.fullPath?o.get(l.fullPath):void 0;if(!d){let t=v(l),r=l.parent;for(;r;){const e=v(r);e&&(t=`${e.replace(/\/$/u,"")}/${t.replace(/^\//u,"")}`),r=r.parent}const i=e(t),n=Object.create(null);for(const e of i)c(e)||(n[e.name]=!0);d={keys:n,tokens:i},o.set(t,d),l.fullPath=t}let p=r(d.tokens,{encode:encodeURIComponent,...i})(s)||"/";if(i.stringifyQueryParams&&s){const e={};for(const[t,r]of Object.entries(s))!(t in d.keys)&&r&&(e[t]=r);const t=i.stringifyQueryParams(e);t&&(p+=t.startsWith("?")?t:`?${t}`)}return p}};const E=/\/\*[\*!]\s+vaadin-dev-mode:start([\s\S]*)vaadin-dev-mode:end\s+\*\*\//i,P=window.Vaadin&&window.Vaadin.Flow&&window.Vaadin.Flow.clients;function R(e,t){if("function"!=typeof e)return;const r=E.exec(e.toString());if(r)try{e=new Function(r[1])}catch(e){console.log("vaadin-development-mode-detector: uncommentAndRun() failed",e)}return e(t)}window.Vaadin=window.Vaadin||{};const O=function(e,t){if(window.Vaadin.developmentMode)return R(e,t)};function S(){
/*! vaadin-dev-mode:start
  (function () {
'use strict';

var _typeof = typeof Symbol === "function" && typeof Symbol.iterator === "symbol" ? function (obj) {
  return typeof obj;
} : function (obj) {
  return obj && typeof Symbol === "function" && obj.constructor === Symbol && obj !== Symbol.prototype ? "symbol" : typeof obj;
};

var classCallCheck = function (instance, Constructor) {
  if (!(instance instanceof Constructor)) {
    throw new TypeError("Cannot call a class as a function");
  }
};

var createClass = function () {
  function defineProperties(target, props) {
    for (var i = 0; i < props.length; i++) {
      var descriptor = props[i];
      descriptor.enumerable = descriptor.enumerable || false;
      descriptor.configurable = true;
      if ("value" in descriptor) descriptor.writable = true;
      Object.defineProperty(target, descriptor.key, descriptor);
    }
  }

  return function (Constructor, protoProps, staticProps) {
    if (protoProps) defineProperties(Constructor.prototype, protoProps);
    if (staticProps) defineProperties(Constructor, staticProps);
    return Constructor;
  };
}();

var getPolymerVersion = function getPolymerVersion() {
  return window.Polymer && window.Polymer.version;
};

var StatisticsGatherer = function () {
  function StatisticsGatherer(logger) {
    classCallCheck(this, StatisticsGatherer);

    this.now = new Date().getTime();
    this.logger = logger;
  }

  createClass(StatisticsGatherer, [{
    key: 'frameworkVersionDetectors',
    value: function frameworkVersionDetectors() {
      return {
        'Flow': function Flow() {
          if (window.Vaadin && window.Vaadin.Flow && window.Vaadin.Flow.clients) {
            var flowVersions = Object.keys(window.Vaadin.Flow.clients).map(function (key) {
              return window.Vaadin.Flow.clients[key];
            }).filter(function (client) {
              return client.getVersionInfo;
            }).map(function (client) {
              return client.getVersionInfo().flow;
            });
            if (flowVersions.length > 0) {
              return flowVersions[0];
            }
          }
        },
        'Vaadin Framework': function VaadinFramework() {
          if (window.vaadin && window.vaadin.clients) {
            var frameworkVersions = Object.values(window.vaadin.clients).filter(function (client) {
              return client.getVersionInfo;
            }).map(function (client) {
              return client.getVersionInfo().vaadinVersion;
            });
            if (frameworkVersions.length > 0) {
              return frameworkVersions[0];
            }
          }
        },
        'AngularJs': function AngularJs() {
          if (window.angular && window.angular.version && window.angular.version) {
            return window.angular.version.full;
          }
        },
        'Angular': function Angular() {
          if (window.ng) {
            var tags = document.querySelectorAll("[ng-version]");
            if (tags.length > 0) {
              return tags[0].getAttribute("ng-version");
            }
            return "Unknown";
          }
        },
        'Backbone.js': function BackboneJs() {
          if (window.Backbone) {
            return window.Backbone.VERSION;
          }
        },
        'React': function React() {
          var reactSelector = '[data-reactroot], [data-reactid]';
          if (!!document.querySelector(reactSelector)) {
            // React does not publish the version by default
            return "unknown";
          }
        },
        'Ember': function Ember() {
          if (window.Em && window.Em.VERSION) {
            return window.Em.VERSION;
          } else if (window.Ember && window.Ember.VERSION) {
            return window.Ember.VERSION;
          }
        },
        'jQuery': function (_jQuery) {
          function jQuery() {
            return _jQuery.apply(this, arguments);
          }

          jQuery.toString = function () {
            return _jQuery.toString();
          };

          return jQuery;
        }(function () {
          if (typeof jQuery === 'function' && jQuery.prototype.jquery !== undefined) {
            return jQuery.prototype.jquery;
          }
        }),
        'Polymer': function Polymer() {
          var version = getPolymerVersion();
          if (version) {
            return version;
          }
        },
        'LitElement': function LitElement() {
          var version = window.litElementVersions && window.litElementVersions[0];
          if (version) {
            return version;
          }
        },
        'LitHtml': function LitHtml() {
          var version = window.litHtmlVersions && window.litHtmlVersions[0];
          if (version) {
            return version;
          }
        },
        'Vue.js': function VueJs() {
          if (window.Vue) {
            return window.Vue.version;
          }
        }
      };
    }
  }, {
    key: 'getUsedVaadinElements',
    value: function getUsedVaadinElements(elements) {
      var version = getPolymerVersion();
      var elementClasses = void 0;
      // NOTE: In case you edit the code here, YOU MUST UPDATE any statistics reporting code in Flow.
      // Check all locations calling the method getEntries() in
      // https://github.com/vaadin/flow/blob/master/flow-server/src/main/java/com/vaadin/flow/internal/UsageStatistics.java#L106
      // Currently it is only used by BootstrapHandler.
      if (version && version.indexOf('2') === 0) {
        // Polymer 2: components classes are stored in window.Vaadin
        elementClasses = Object.keys(window.Vaadin).map(function (c) {
          return window.Vaadin[c];
        }).filter(function (c) {
          return c.is;
        });
      } else {
        // Polymer 3: components classes are stored in window.Vaadin.registrations
        elementClasses = window.Vaadin.registrations || [];
      }
      elementClasses.forEach(function (klass) {
        var version = klass.version ? klass.version : "0.0.0";
        elements[klass.is] = { version: version };
      });
    }
  }, {
    key: 'getUsedVaadinThemes',
    value: function getUsedVaadinThemes(themes) {
      ['Lumo', 'Material'].forEach(function (themeName) {
        var theme;
        var version = getPolymerVersion();
        if (version && version.indexOf('2') === 0) {
          // Polymer 2: themes are stored in window.Vaadin
          theme = window.Vaadin[themeName];
        } else {
          // Polymer 3: themes are stored in custom element registry
          theme = customElements.get('vaadin-' + themeName.toLowerCase() + '-styles');
        }
        if (theme && theme.version) {
          themes[themeName] = { version: theme.version };
        }
      });
    }
  }, {
    key: 'getFrameworks',
    value: function getFrameworks(frameworks) {
      var detectors = this.frameworkVersionDetectors();
      Object.keys(detectors).forEach(function (framework) {
        var detector = detectors[framework];
        try {
          var version = detector();
          if (version) {
            frameworks[framework] = { version: version };
          }
        } catch (e) {}
      });
    }
  }, {
    key: 'gather',
    value: function gather(storage) {
      var storedStats = storage.read();
      var gatheredStats = {};
      var types = ["elements", "frameworks", "themes"];

      types.forEach(function (type) {
        gatheredStats[type] = {};
        if (!storedStats[type]) {
          storedStats[type] = {};
        }
      });

      var previousStats = JSON.stringify(storedStats);

      this.getUsedVaadinElements(gatheredStats.elements);
      this.getFrameworks(gatheredStats.frameworks);
      this.getUsedVaadinThemes(gatheredStats.themes);

      var now = this.now;
      types.forEach(function (type) {
        var keys = Object.keys(gatheredStats[type]);
        keys.forEach(function (key) {
          if (!storedStats[type][key] || _typeof(storedStats[type][key]) != _typeof({})) {
            storedStats[type][key] = { firstUsed: now };
          }
          // Discards any previously logged version number
          storedStats[type][key].version = gatheredStats[type][key].version;
          storedStats[type][key].lastUsed = now;
        });
      });

      var newStats = JSON.stringify(storedStats);
      storage.write(newStats);
      if (newStats != previousStats && Object.keys(storedStats).length > 0) {
        this.logger.debug("New stats: " + newStats);
      }
    }
  }]);
  return StatisticsGatherer;
}();

var StatisticsStorage = function () {
  function StatisticsStorage(key) {
    classCallCheck(this, StatisticsStorage);

    this.key = key;
  }

  createClass(StatisticsStorage, [{
    key: 'read',
    value: function read() {
      var localStorageStatsString = localStorage.getItem(this.key);
      try {
        return JSON.parse(localStorageStatsString ? localStorageStatsString : '{}');
      } catch (e) {
        return {};
      }
    }
  }, {
    key: 'write',
    value: function write(data) {
      localStorage.setItem(this.key, data);
    }
  }, {
    key: 'clear',
    value: function clear() {
      localStorage.removeItem(this.key);
    }
  }, {
    key: 'isEmpty',
    value: function isEmpty() {
      var storedStats = this.read();
      var empty = true;
      Object.keys(storedStats).forEach(function (key) {
        if (Object.keys(storedStats[key]).length > 0) {
          empty = false;
        }
      });

      return empty;
    }
  }]);
  return StatisticsStorage;
}();

var StatisticsSender = function () {
  function StatisticsSender(url, logger) {
    classCallCheck(this, StatisticsSender);

    this.url = url;
    this.logger = logger;
  }

  createClass(StatisticsSender, [{
    key: 'send',
    value: function send(data, errorHandler) {
      var logger = this.logger;

      if (navigator.onLine === false) {
        logger.debug("Offline, can't send");
        errorHandler();
        return;
      }
      logger.debug("Sending data to " + this.url);

      var req = new XMLHttpRequest();
      req.withCredentials = true;
      req.addEventListener("load", function () {
        // Stats sent, nothing more to do
        logger.debug("Response: " + req.responseText);
      });
      req.addEventListener("error", function () {
        logger.debug("Send failed");
        errorHandler();
      });
      req.addEventListener("abort", function () {
        logger.debug("Send aborted");
        errorHandler();
      });
      req.open("POST", this.url);
      req.setRequestHeader("Content-Type", "application/json");
      req.send(data);
    }
  }]);
  return StatisticsSender;
}();

var StatisticsLogger = function () {
  function StatisticsLogger(id) {
    classCallCheck(this, StatisticsLogger);

    this.id = id;
  }

  createClass(StatisticsLogger, [{
    key: '_isDebug',
    value: function _isDebug() {
      return localStorage.getItem("vaadin." + this.id + ".debug");
    }
  }, {
    key: 'debug',
    value: function debug(msg) {
      if (this._isDebug()) {
        console.info(this.id + ": " + msg);
      }
    }
  }]);
  return StatisticsLogger;
}();

var UsageStatistics = function () {
  function UsageStatistics() {
    classCallCheck(this, UsageStatistics);

    this.now = new Date();
    this.timeNow = this.now.getTime();
    this.gatherDelay = 10; // Delay between loading this file and gathering stats
    this.initialDelay = 24 * 60 * 60;

    this.logger = new StatisticsLogger("statistics");
    this.storage = new StatisticsStorage("vaadin.statistics.basket");
    this.gatherer = new StatisticsGatherer(this.logger);
    this.sender = new StatisticsSender("https://tools.vaadin.com/usage-stats/submit", this.logger);
  }

  createClass(UsageStatistics, [{
    key: 'maybeGatherAndSend',
    value: function maybeGatherAndSend() {
      var _this = this;

      if (localStorage.getItem(UsageStatistics.optOutKey)) {
        return;
      }
      this.gatherer.gather(this.storage);
      setTimeout(function () {
        _this.maybeSend();
      }, this.gatherDelay * 1000);
    }
  }, {
    key: 'lottery',
    value: function lottery() {
      return true;
    }
  }, {
    key: 'currentMonth',
    value: function currentMonth() {
      return this.now.getYear() * 12 + this.now.getMonth();
    }
  }, {
    key: 'maybeSend',
    value: function maybeSend() {
      var firstUse = Number(localStorage.getItem(UsageStatistics.firstUseKey));
      var monthProcessed = Number(localStorage.getItem(UsageStatistics.monthProcessedKey));

      if (!firstUse) {
        // Use a grace period to avoid interfering with tests, incognito mode etc
        firstUse = this.timeNow;
        localStorage.setItem(UsageStatistics.firstUseKey, firstUse);
      }

      if (this.timeNow < firstUse + this.initialDelay * 1000) {
        this.logger.debug("No statistics will be sent until the initial delay of " + this.initialDelay + "s has passed");
        return;
      }
      if (this.currentMonth() <= monthProcessed) {
        this.logger.debug("This month has already been processed");
        return;
      }
      localStorage.setItem(UsageStatistics.monthProcessedKey, this.currentMonth());
      // Use random sampling
      if (this.lottery()) {
        this.logger.debug("Congratulations, we have a winner!");
      } else {
        this.logger.debug("Sorry, no stats from you this time");
        return;
      }

      this.send();
    }
  }, {
    key: 'send',
    value: function send() {
      // Ensure we have the latest data
      this.gatherer.gather(this.storage);

      // Read, send and clean up
      var data = this.storage.read();
      data["firstUse"] = Number(localStorage.getItem(UsageStatistics.firstUseKey));
      data["usageStatisticsVersion"] = UsageStatistics.version;
      var info = 'This request contains usage statistics gathered from the application running in development mode. \n\nStatistics gathering is automatically disabled and excluded from production builds.\n\nFor details and to opt-out, see https://github.com/vaadin/vaadin-usage-statistics.\n\n\n\n';
      var self = this;
      this.sender.send(info + JSON.stringify(data), function () {
        // Revert the 'month processed' flag
        localStorage.setItem(UsageStatistics.monthProcessedKey, self.currentMonth() - 1);
      });
    }
  }], [{
    key: 'version',
    get: function get$1() {
      return '2.1.2';
    }
  }, {
    key: 'firstUseKey',
    get: function get$1() {
      return 'vaadin.statistics.firstuse';
    }
  }, {
    key: 'monthProcessedKey',
    get: function get$1() {
      return 'vaadin.statistics.monthProcessed';
    }
  }, {
    key: 'optOutKey',
    get: function get$1() {
      return 'vaadin.statistics.optout';
    }
  }]);
  return UsageStatistics;
}();

try {
  window.Vaadin = window.Vaadin || {};
  window.Vaadin.usageStatsChecker = window.Vaadin.usageStatsChecker || new UsageStatistics();
  window.Vaadin.usageStatsChecker.maybeGatherAndSend();
} catch (e) {
  // Intentionally ignored as this is not a problem in the app being developed
}

}());

  vaadin-dev-mode:end **/}void 0===window.Vaadin.developmentMode&&(window.Vaadin.developmentMode=function(){try{return!!localStorage.getItem("vaadin.developmentmode.force")||["localhost","127.0.0.1"].indexOf(window.location.hostname)>=0&&(P?!(P&&Object.keys(P).map((e=>P[e])).filter((e=>e.productionMode)).length>0):!R((function(){return!0})))}catch(e){return!1}}());!function(e,t=(window.Vaadin??={})){t.registrations??=[],t.registrations.push({is:"@vaadin/router",version:"2.0.0"})}(),O(S);var z=async function(e,t){return e.classList.add(t),await new Promise((r=>{if((e=>{const t=getComputedStyle(e).getPropertyValue("animation-name");return t&&"none"!==t})(e)){const i=e.getBoundingClientRect(),o=`height: ${i.bottom-i.top}px; width: ${i.right-i.left}px`;e.setAttribute("style",`position: absolute; ${o}`),((e,t)=>{const r=()=>{e.removeEventListener("animationend",r),t()};e.addEventListener("animationend",r)})(e,(()=>{e.classList.remove(t),e.removeAttribute("style"),r()}))}else e.classList.remove(t),r()}))};function T(e){if(!e||!c(e.path))throw new Error(p('Expected route config to be an object with a "path" string property, or an array of such objects'));if(!(l(e.action)||Array.isArray(e.children)||l(e.children)||c(e.component)||c(e.redirect)))throw new Error(p(`Expected route config "${e.path}" to include either "component, redirect" or "action" function but none found.`));e.redirect&&["bundle","component"].forEach((t=>{t in e&&console.warn(p(`Route config "${String(e.path)}" has both "redirect" and "${t}" properties, and "redirect" will always override the latter. Did you mean to only use "${t}"?`))}))}function M(e){d(e).forEach((e=>T(e)))}function L(e,t){const r=t.__effectiveBaseUrl;return r?new URL(e.replace(/^\//u,""),r).pathname:e}function j(e){return e.map((e=>e.path)).reduce(((e,t)=>t.length?`${e.replace(/\/$/u,"")}/${t.replace(/^\//u,"")}`:e),"")}function B({chain:e=[],hash:r="",params:i={},pathname:o="",redirectFrom:n,resolver:a,search:s=""},l){const c=e.map((e=>e.route));return{baseUrl:a?.baseUrl??"",getUrl:(r={})=>a?L(t(function(e){return j(e.map((e=>e.route)))}(e))({...i,...r}),a):"",hash:r,params:i,pathname:o,redirectFrom:n,route:l??(Array.isArray(c)?c.at(-1):void 0)??null,routes:c,search:s,searchParams:new URLSearchParams(s)}}function I(e,t){const r={...e.params};return{redirect:{from:e.pathname,params:r,pathname:t}}}function N(e,t,...r){if("function"==typeof e)return e.apply(t,r)}function U(e,t,...r){return i=>i&&s(i)&&("cancel"in i||"redirect"in i)?i:N(t?.[e],t,...r)}function D(e,t){return!window.dispatchEvent(new CustomEvent(`vaadin-router-${e}`,{cancelable:"go"===e,detail:t}))}function H(e){if(e instanceof Element)return e.nodeName.toLowerCase()}function V(e){if(e.defaultPrevented)return;if(0!==e.button)return;if(e.shiftKey||e.ctrlKey||e.altKey||e.metaKey)return;let t=e.target;const r=e instanceof MouseEvent?e.composedPath():e.path??[];for(let e=0;e<r.length;e++){const i=r[e];if("nodeName"in i&&"a"===i.nodeName.toLowerCase()){t=i;break}}for(;t&&t instanceof Node&&"a"!==H(t);)t=t.parentNode;if(!t||"a"!==H(t))return;const i=t;if(i.target&&"_self"!==i.target.toLowerCase())return;if(i.hasAttribute("download"))return;if(i.hasAttribute("router-ignore"))return;if(i.pathname===window.location.pathname&&""!==i.hash)return;const o=i.origin||function(e){const{port:t,protocol:r}=e;return`${r}//${"http:"===r&&"80"===t||"https:"===r&&"443"===t?e.hostname:e.host}`}(i);if(o!==window.location.origin)return;const{hash:n,pathname:a,search:s}=i;D("go",{hash:n,pathname:a,search:s})&&e instanceof MouseEvent&&(e.preventDefault(),"click"===e.type&&window.scrollTo(0,0))}function K(e){if("vaadin-router-ignore"===e.state)return;const{hash:t,pathname:r,search:i}=window.location;D("go",{hash:t,pathname:r,search:i})}let F=[];const W={CLICK:{activate(){window.document.addEventListener("click",V)},inactivate(){window.document.removeEventListener("click",V)}},POPSTATE:{activate(){window.addEventListener("popstate",K)},inactivate(){window.removeEventListener("popstate",K)}}};function Y(e=[]){F.forEach((e=>e.inactivate())),e.forEach((e=>e.activate())),F=e}function G(){return{cancel:!0}}const q={__renderId:-1,params:{},route:{__synthetic:!0,children:[],path:"",action(){}},pathname:"",next:async()=>u};class J extends k{location=B({resolver:this});ready=Promise.resolve(this.location);#r=new WeakSet;#i=new WeakSet;#o=this.#n.bind(this);#a=0;#s;__previousContext;#l;#c=null;#d=null;constructor(e,t){const r=document.head.querySelector("base"),i=r?.getAttribute("href");super([],{baseUrl:i?new URL(i,document.URL).href.replace(/[^/]*$/u,""):void 0,...t,resolveRoute:async e=>await this.#p(e)}),Y(Object.values(W)),this.setOutlet(e),this.subscribe()}async#p(e){const{route:t}=e;if(l(t.children)){let r=await t.children(function({next:e,...t}){return t}(e));l(t.children)||({children:r}=t),function(e,t){if(!Array.isArray(e)&&!s(e))throw new Error(p(`Incorrect "children" value for the route ${String(t.path)}: expected array or object, but got ${String(e)}`));const r=d(e);r.forEach((e=>T(e))),t.__children=r}(r,t)}const r={component:e=>{const t=document.createElement(e);return this.#i.add(t),t},prevent:G,redirect:t=>I(e,t)};return await Promise.resolve().then((async()=>{if(this.#h(e))return await N(t.action,t,e,r)})).then((e=>null==e||"object"!=typeof e&&"symbol"!=typeof e||!(e instanceof HTMLElement||e===u||s(e)&&"redirect"in e)?c(t.redirect)?r.redirect(t.redirect):void 0:e)).then((e=>null!=e?e:c(t.component)?r.component(t.component):void 0))}setOutlet(e){e&&this.#u(e),this.#s=e}getOutlet(){return this.#s}async setRoutes(e,t=!1){return this.__previousContext=void 0,this.#l=void 0,M(e),super.setRoutes(e),t||this.#n(),await this.ready}addRoutes(e){return M(e),super.addRoutes(e)}async render(e,t=!1){this.#a+=1;const r=this.#a,i={...q,...c(e)?{hash:"",search:"",pathname:e}:e,__renderId:r};return this.ready=this.#m(i,t),await this.ready}async#m(e,t){const{__renderId:r}=e;try{const i=await this.resolve(e),o=await this.#f(i);if(!this.#h(o))return this.location;const n=this.__previousContext;if(o===n)return this.#v(n,!0),this.location;if(this.location=B(o),t&&this.#v(o,1===r),D("location-changed",{router:this,location:this.location}),o.__skipAttach)return this.#g(o,n),this.__previousContext=o,this.location;this.#y(o,n);const a=this.#b(o);if(this.#w(o),this.#x(o,n),await a,this.#h(o))return this.#$(),this.__previousContext=o,this.location}catch(i){if(r===this.#a){t&&this.#v(this.context);for(const e of this.#s?.children??[])e.remove();throw this.location=B(Object.assign(e,{resolver:this})),D("error",{router:this,error:i,...e}),i}}return this.location}async#f(e,t=e){const r=await this.#k(t),i=r!==t?r:e,o=L(j(r.chain??[]),this)===r.pathname,n=async(e,t=e.route,r)=>{const i=await e.next(!1,t,r);return null===i||i===u?o?e:null!=t.parent?await n(e,t.parent,i):i:i},a=await n(r);if(null==a||a===u)throw m(i);return a!==r?await this.#f(i,a):await this.#_(r)}async#k(e){const{result:t}=e;if(t instanceof HTMLElement)return function(e,t){if(t.location=B(e),e.chain){const r=e.chain.map((e=>e.route)).indexOf(e.route);e.chain[r].element=t}}(e,t),e;if(t&&"redirect"in t){const r=await this.#A(t.redirect,e.__redirectCount,e.__renderId);return await this.#k(r)}throw t instanceof Error?t:new Error(p(`Invalid route resolution result for path "${e.pathname}". Expected redirect object or HTML element, but got: "${function(e){if("object"!=typeof e)return String(e);const[t="Unknown"]=/ (.*)\]$/u.exec(String(e))??[];return"Object"===t||"Array"===t?`${t} ${JSON.stringify(e)}`:t}(t)}". Double check the action return value for the route.`))}async#_(e){return await this.#C(e).then((async t=>t===this.__previousContext||t===e?t:await this.#f(t)))}async#C(e){const t=this.__previousContext??{},r=t.chain??[],i=e.chain??[];let o=Promise.resolve(void 0);const n=t=>I(e,t);if(e.__divergedChainIndex=0,e.__skipAttach=!1,r.length){for(let t=0;t<Math.min(r.length,i.length)&&(r[t].route===i[t].route&&(r[t].path===i[t].path||r[t].element===i[t].element)&&this.#E(r[t].element,i[t].element));e.__divergedChainIndex++,t++);if(e.__skipAttach=i.length===r.length&&e.__divergedChainIndex===i.length&&this.#E(e.result,t.result),e.__skipAttach){for(let t=i.length-1;t>=0;t--)o=this.#P(o,e,{prevent:G},r[t]);for(let t=0;t<i.length;t++)o=this.#R(o,e,{prevent:G,redirect:n},i[t]),r[t].element.location=B(e,r[t].route)}else for(let t=r.length-1;t>=e.__divergedChainIndex;t--)o=this.#P(o,e,{prevent:G},r[t])}if(!e.__skipAttach)for(let t=0;t<i.length;t++)t<e.__divergedChainIndex?t<r.length&&r[t].element&&(r[t].element.location=B(e,r[t].route)):(o=this.#R(o,e,{prevent:G,redirect:n},i[t]),i[t].element&&(i[t].element.location=B(e,i[t].route)));return await o.then((async t=>{if(t&&s(t)){if("cancel"in t&&this.__previousContext)return this.__previousContext.__renderId=e.__renderId,this.__previousContext;if("redirect"in t)return await this.#A(t.redirect,e.__redirectCount,e.__renderId)}return e}))}async#P(e,t,r,i){const o=B(t);let n=await e;if(this.#h(t)){n=U("onBeforeLeave",i.element,o,r,this)(n)}if(!s(n)||!("redirect"in n))return n}async#R(e,t,r,i){const o=B(t,i.route),n=await e;if(this.#h(t)){return U("onBeforeEnter",i.element,o,r,this)(n)}}#E(e,t){return e instanceof Element&&t instanceof Element&&(this.#i.has(e)&&this.#i.has(t)?e.localName===t.localName:e===t)}#h(e){return e.__renderId===this.#a}async#A(e,t=0,r=0){if(t>256)throw new Error(p(`Too many redirects when rendering ${e.from}`));return await this.resolve({...q,pathname:this.urlForPath(e.pathname,e.params),redirectFrom:e.from,__redirectCount:t+1,__renderId:r})}#u(e=this.#s){if(!(e instanceof Element||e instanceof DocumentFragment))throw new TypeError(p(`Expected router outlet to be a valid DOM Element | DocumentFragment (but got ${e})`))}#v({pathname:e,search:t="",hash:r=""},i){if(window.location.pathname!==e||window.location.search!==t||window.location.hash!==r){const o=i?"replaceState":"pushState";window.history[o](null,document.title,e+t+r),window.dispatchEvent(new PopStateEvent("popstate",{state:"vaadin-router-ignore"}))}}#g(e,t){let r=this.#s;for(let i=0;i<(e.__divergedChainIndex??0);i++){const o=t?.chain?.[i].element;if(o){if(o.parentNode!==r)break;e.chain[i].element=o,r=o}}return r}#y(e,t){this.#u(),this.#O();const r=this.#g(e,t);this.#c=[],this.#d=Array.from(r?.children??[]).filter((t=>this.#r.has(t)&&t!==e.result));let i=r;for(let t=e.__divergedChainIndex??0;t<(e.chain?.length??0);t++){const o=e.chain[t].element;o&&(i?.appendChild(o),this.#r.add(o),i===r&&this.#c.push(o),i=o)}}#$(){if(this.#d)for(const e of this.#d)e.remove();this.#d=null,this.#c=null}#O(){if(this.#d&&this.#c){for(const e of this.#c)e.remove();this.#d=null,this.#c=null}}#x(e,t){if(t?.chain&&null!=e.__divergedChainIndex)for(let r=t.chain.length-1;r>=e.__divergedChainIndex&&this.#h(e);r--){const i=t.chain[r].element;if(i)try{const t=B(e);N(i.onAfterLeave,i,t,{},this)}finally{if(this.#d?.includes(i))for(const e of i.children)e.remove()}}}#w(e){if(e.chain&&null!=e.__divergedChainIndex)for(let t=e.__divergedChainIndex;t<e.chain.length&&this.#h(e);t++){const r=e.chain[t].element;if(r){const i=B(e,e.chain[t].route);N(r.onAfterEnter,r,i,{},this)}}}async#b(e){const t=this.#d?.[0],r=this.#c?.[0],i=[],{chain:o=[]}=e;let n;for(let e=o.length-1;e>=0;e--)if(o[e].route.animate){n=o[e].route.animate;break}if(t&&r&&n){const e=s(n)&&n.leave?n.leave:"leaving",o=s(n)&&n.enter?n.enter:"entering";i.push(z(t,e)),i.push(z(r,o))}return await Promise.all(i),e}subscribe(){window.addEventListener("vaadin-router-go",this.#o)}unsubscribe(){window.removeEventListener("vaadin-router-go",this.#o)}#n(e){const{pathname:t,search:r,hash:i}=e instanceof CustomEvent?e.detail:window.location;c(this.__normalizePathname(t))&&(e?.preventDefault&&e.preventDefault(),this.render({pathname:t,search:r,hash:i},!0))}static setTriggers(...e){Y(e)}urlForName(e,t){return this.#l||(this.#l=C(this,{cacheKeyProvider:e=>"component"in e&&"string"==typeof e.component?e.component:void 0})),L(this.#l(e,t??void 0),this)}urlForPath(e,r){return L(t(e)(r??void 0),this)}static go(e){const{pathname:t,search:r,hash:i}=c(e)?new URL(e,"http://a"):e;return D("go",{pathname:t,search:r,hash:i})}}
/**
 * @license
 * Copyright 2019 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */const Q=globalThis,Z=Q.ShadowRoot&&(void 0===Q.ShadyCSS||Q.ShadyCSS.nativeShadow)&&"adoptedStyleSheets"in Document.prototype&&"replace"in CSSStyleSheet.prototype,X=Symbol(),ee=new WeakMap;let te=class{constructor(e,t,r){if(this._$cssResult$=!0,r!==X)throw Error("CSSResult is not constructable. Use `unsafeCSS` or `css` instead.");this.cssText=e,this.t=t}get styleSheet(){let e=this.o;const t=this.t;if(Z&&void 0===e){const r=void 0!==t&&1===t.length;r&&(e=ee.get(t)),void 0===e&&((this.o=e=new CSSStyleSheet).replaceSync(this.cssText),r&&ee.set(t,e))}return e}toString(){return this.cssText}};const re=(e,...t)=>{const r=1===e.length?e[0]:t.reduce(((t,r,i)=>t+(e=>{if(!0===e._$cssResult$)return e.cssText;if("number"==typeof e)return e;throw Error("Value passed to 'css' function must be a 'css' function result: "+e+". Use 'unsafeCSS' to pass non-literal values, but take care to ensure page security.")})(r)+e[i+1]),e[0]);return new te(r,e,X)},ie=Z?e=>e:e=>e instanceof CSSStyleSheet?(e=>{let t="";for(const r of e.cssRules)t+=r.cssText;return(e=>new te("string"==typeof e?e:e+"",void 0,X))(t)})(e):e
/**
 * @license
 * Copyright 2017 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */,{is:oe,defineProperty:ne,getOwnPropertyDescriptor:ae,getOwnPropertyNames:se,getOwnPropertySymbols:le,getPrototypeOf:ce}=Object,de=globalThis,pe=de.trustedTypes,he=pe?pe.emptyScript:"",ue=de.reactiveElementPolyfillSupport,me=(e,t)=>e,fe={toAttribute(e,t){switch(t){case Boolean:e=e?he:null;break;case Object:case Array:e=null==e?e:JSON.stringify(e)}return e},fromAttribute(e,t){let r=e;switch(t){case Boolean:r=null!==e;break;case Number:r=null===e?null:Number(e);break;case Object:case Array:try{r=JSON.parse(e)}catch(e){r=null}}return r}},ve=(e,t)=>!oe(e,t),ge={attribute:!0,type:String,converter:fe,reflect:!1,hasChanged:ve};Symbol.metadata??=Symbol("metadata"),de.litPropertyMetadata??=new WeakMap;class ye extends HTMLElement{static addInitializer(e){this._$Ei(),(this.l??=[]).push(e)}static get observedAttributes(){return this.finalize(),this._$Eh&&[...this._$Eh.keys()]}static createProperty(e,t=ge){if(t.state&&(t.attribute=!1),this._$Ei(),this.elementProperties.set(e,t),!t.noAccessor){const r=Symbol(),i=this.getPropertyDescriptor(e,r,t);void 0!==i&&ne(this.prototype,e,i)}}static getPropertyDescriptor(e,t,r){const{get:i,set:o}=ae(this.prototype,e)??{get(){return this[t]},set(e){this[t]=e}};return{get(){return i?.call(this)},set(t){const n=i?.call(this);o.call(this,t),this.requestUpdate(e,n,r)},configurable:!0,enumerable:!0}}static getPropertyOptions(e){return this.elementProperties.get(e)??ge}static _$Ei(){if(this.hasOwnProperty(me("elementProperties")))return;const e=ce(this);e.finalize(),void 0!==e.l&&(this.l=[...e.l]),this.elementProperties=new Map(e.elementProperties)}static finalize(){if(this.hasOwnProperty(me("finalized")))return;if(this.finalized=!0,this._$Ei(),this.hasOwnProperty(me("properties"))){const e=this.properties,t=[...se(e),...le(e)];for(const r of t)this.createProperty(r,e[r])}const e=this[Symbol.metadata];if(null!==e){const t=litPropertyMetadata.get(e);if(void 0!==t)for(const[e,r]of t)this.elementProperties.set(e,r)}this._$Eh=new Map;for(const[e,t]of this.elementProperties){const r=this._$Eu(e,t);void 0!==r&&this._$Eh.set(r,e)}this.elementStyles=this.finalizeStyles(this.styles)}static finalizeStyles(e){const t=[];if(Array.isArray(e)){const r=new Set(e.flat(1/0).reverse());for(const e of r)t.unshift(ie(e))}else void 0!==e&&t.push(ie(e));return t}static _$Eu(e,t){const r=t.attribute;return!1===r?void 0:"string"==typeof r?r:"string"==typeof e?e.toLowerCase():void 0}constructor(){super(),this._$Ep=void 0,this.isUpdatePending=!1,this.hasUpdated=!1,this._$Em=null,this._$Ev()}_$Ev(){this._$ES=new Promise((e=>this.enableUpdating=e)),this._$AL=new Map,this._$E_(),this.requestUpdate(),this.constructor.l?.forEach((e=>e(this)))}addController(e){(this._$EO??=new Set).add(e),void 0!==this.renderRoot&&this.isConnected&&e.hostConnected?.()}removeController(e){this._$EO?.delete(e)}_$E_(){const e=new Map,t=this.constructor.elementProperties;for(const r of t.keys())this.hasOwnProperty(r)&&(e.set(r,this[r]),delete this[r]);e.size>0&&(this._$Ep=e)}createRenderRoot(){const e=this.shadowRoot??this.attachShadow(this.constructor.shadowRootOptions);return((e,t)=>{if(Z)e.adoptedStyleSheets=t.map((e=>e instanceof CSSStyleSheet?e:e.styleSheet));else for(const r of t){const t=document.createElement("style"),i=Q.litNonce;void 0!==i&&t.setAttribute("nonce",i),t.textContent=r.cssText,e.appendChild(t)}})(e,this.constructor.elementStyles),e}connectedCallback(){this.renderRoot??=this.createRenderRoot(),this.enableUpdating(!0),this._$EO?.forEach((e=>e.hostConnected?.()))}enableUpdating(e){}disconnectedCallback(){this._$EO?.forEach((e=>e.hostDisconnected?.()))}attributeChangedCallback(e,t,r){this._$AK(e,r)}_$EC(e,t){const r=this.constructor.elementProperties.get(e),i=this.constructor._$Eu(e,r);if(void 0!==i&&!0===r.reflect){const o=(void 0!==r.converter?.toAttribute?r.converter:fe).toAttribute(t,r.type);this._$Em=e,null==o?this.removeAttribute(i):this.setAttribute(i,o),this._$Em=null}}_$AK(e,t){const r=this.constructor,i=r._$Eh.get(e);if(void 0!==i&&this._$Em!==i){const e=r.getPropertyOptions(i),o="function"==typeof e.converter?{fromAttribute:e.converter}:void 0!==e.converter?.fromAttribute?e.converter:fe;this._$Em=i,this[i]=o.fromAttribute(t,e.type),this._$Em=null}}requestUpdate(e,t,r){if(void 0!==e){if(r??=this.constructor.getPropertyOptions(e),!(r.hasChanged??ve)(this[e],t))return;this.P(e,t,r)}!1===this.isUpdatePending&&(this._$ES=this._$ET())}P(e,t,r){this._$AL.has(e)||this._$AL.set(e,t),!0===r.reflect&&this._$Em!==e&&(this._$Ej??=new Set).add(e)}async _$ET(){this.isUpdatePending=!0;try{await this._$ES}catch(e){Promise.reject(e)}const e=this.scheduleUpdate();return null!=e&&await e,!this.isUpdatePending}scheduleUpdate(){return this.performUpdate()}performUpdate(){if(!this.isUpdatePending)return;if(!this.hasUpdated){if(this.renderRoot??=this.createRenderRoot(),this._$Ep){for(const[e,t]of this._$Ep)this[e]=t;this._$Ep=void 0}const e=this.constructor.elementProperties;if(e.size>0)for(const[t,r]of e)!0!==r.wrapped||this._$AL.has(t)||void 0===this[t]||this.P(t,this[t],r)}let e=!1;const t=this._$AL;try{e=this.shouldUpdate(t),e?(this.willUpdate(t),this._$EO?.forEach((e=>e.hostUpdate?.())),this.update(t)):this._$EU()}catch(t){throw e=!1,this._$EU(),t}e&&this._$AE(t)}willUpdate(e){}_$AE(e){this._$EO?.forEach((e=>e.hostUpdated?.())),this.hasUpdated||(this.hasUpdated=!0,this.firstUpdated(e)),this.updated(e)}_$EU(){this._$AL=new Map,this.isUpdatePending=!1}get updateComplete(){return this.getUpdateComplete()}getUpdateComplete(){return this._$ES}shouldUpdate(e){return!0}update(e){this._$Ej&&=this._$Ej.forEach((e=>this._$EC(e,this[e]))),this._$EU()}updated(e){}firstUpdated(e){}}ye.elementStyles=[],ye.shadowRootOptions={mode:"open"},ye[me("elementProperties")]=new Map,ye[me("finalized")]=new Map,ue?.({ReactiveElement:ye}),(de.reactiveElementVersions??=[]).push("2.0.4");
/**
 * @license
 * Copyright 2017 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */
const be=globalThis,we=be.trustedTypes,xe=we?we.createPolicy("lit-html",{createHTML:e=>e}):void 0,$e="$lit$",ke=`lit$${Math.random().toFixed(9).slice(2)}$`,_e="?"+ke,Ae=`<${_e}>`,Ce=document,Ee=()=>Ce.createComment(""),Pe=e=>null===e||"object"!=typeof e&&"function"!=typeof e,Re=Array.isArray,Oe="[ \t\n\f\r]",Se=/<(?:(!--|\/[^a-zA-Z])|(\/?[a-zA-Z][^>\s]*)|(\/?$))/g,ze=/-->/g,Te=/>/g,Me=RegExp(`>|${Oe}(?:([^\\s"'>=/]+)(${Oe}*=${Oe}*(?:[^ \t\n\f\r"'\`<>=]|("|')|))|$)`,"g"),Le=/'/g,je=/"/g,Be=/^(?:script|style|textarea|title)$/i,Ie=(e=>(t,...r)=>({_$litType$:e,strings:t,values:r}))(1),Ne=Symbol.for("lit-noChange"),Ue=Symbol.for("lit-nothing"),De=new WeakMap,He=Ce.createTreeWalker(Ce,129);function Ve(e,t){if(!Re(e)||!e.hasOwnProperty("raw"))throw Error("invalid template strings array");return void 0!==xe?xe.createHTML(t):t}class Ke{constructor({strings:e,_$litType$:t},r){let i;this.parts=[];let o=0,n=0;const a=e.length-1,s=this.parts,[l,c]=((e,t)=>{const r=e.length-1,i=[];let o,n=2===t?"<svg>":3===t?"<math>":"",a=Se;for(let t=0;t<r;t++){const r=e[t];let s,l,c=-1,d=0;for(;d<r.length&&(a.lastIndex=d,l=a.exec(r),null!==l);)d=a.lastIndex,a===Se?"!--"===l[1]?a=ze:void 0!==l[1]?a=Te:void 0!==l[2]?(Be.test(l[2])&&(o=RegExp("</"+l[2],"g")),a=Me):void 0!==l[3]&&(a=Me):a===Me?">"===l[0]?(a=o??Se,c=-1):void 0===l[1]?c=-2:(c=a.lastIndex-l[2].length,s=l[1],a=void 0===l[3]?Me:'"'===l[3]?je:Le):a===je||a===Le?a=Me:a===ze||a===Te?a=Se:(a=Me,o=void 0);const p=a===Me&&e[t+1].startsWith("/>")?" ":"";n+=a===Se?r+Ae:c>=0?(i.push(s),r.slice(0,c)+$e+r.slice(c)+ke+p):r+ke+(-2===c?t:p)}return[Ve(e,n+(e[r]||"<?>")+(2===t?"</svg>":3===t?"</math>":"")),i]})(e,t);if(this.el=Ke.createElement(l,r),He.currentNode=this.el.content,2===t||3===t){const e=this.el.content.firstChild;e.replaceWith(...e.childNodes)}for(;null!==(i=He.nextNode())&&s.length<a;){if(1===i.nodeType){if(i.hasAttributes())for(const e of i.getAttributeNames())if(e.endsWith($e)){const t=c[n++],r=i.getAttribute(e).split(ke),a=/([.?@])?(.*)/.exec(t);s.push({type:1,index:o,name:a[2],strings:r,ctor:"."===a[1]?qe:"?"===a[1]?Je:"@"===a[1]?Qe:Ge}),i.removeAttribute(e)}else e.startsWith(ke)&&(s.push({type:6,index:o}),i.removeAttribute(e));if(Be.test(i.tagName)){const e=i.textContent.split(ke),t=e.length-1;if(t>0){i.textContent=we?we.emptyScript:"";for(let r=0;r<t;r++)i.append(e[r],Ee()),He.nextNode(),s.push({type:2,index:++o});i.append(e[t],Ee())}}}else if(8===i.nodeType)if(i.data===_e)s.push({type:2,index:o});else{let e=-1;for(;-1!==(e=i.data.indexOf(ke,e+1));)s.push({type:7,index:o}),e+=ke.length-1}o++}}static createElement(e,t){const r=Ce.createElement("template");return r.innerHTML=e,r}}function Fe(e,t,r=e,i){if(t===Ne)return t;let o=void 0!==i?r._$Co?.[i]:r._$Cl;const n=Pe(t)?void 0:t._$litDirective$;return o?.constructor!==n&&(o?._$AO?.(!1),void 0===n?o=void 0:(o=new n(e),o._$AT(e,r,i)),void 0!==i?(r._$Co??=[])[i]=o:r._$Cl=o),void 0!==o&&(t=Fe(e,o._$AS(e,t.values),o,i)),t}let We=class{constructor(e,t){this._$AV=[],this._$AN=void 0,this._$AD=e,this._$AM=t}get parentNode(){return this._$AM.parentNode}get _$AU(){return this._$AM._$AU}u(e){const{el:{content:t},parts:r}=this._$AD,i=(e?.creationScope??Ce).importNode(t,!0);He.currentNode=i;let o=He.nextNode(),n=0,a=0,s=r[0];for(;void 0!==s;){if(n===s.index){let t;2===s.type?t=new Ye(o,o.nextSibling,this,e):1===s.type?t=new s.ctor(o,s.name,s.strings,this,e):6===s.type&&(t=new Ze(o,this,e)),this._$AV.push(t),s=r[++a]}n!==s?.index&&(o=He.nextNode(),n++)}return He.currentNode=Ce,i}p(e){let t=0;for(const r of this._$AV)void 0!==r&&(void 0!==r.strings?(r._$AI(e,r,t),t+=r.strings.length-2):r._$AI(e[t])),t++}};class Ye{get _$AU(){return this._$AM?._$AU??this._$Cv}constructor(e,t,r,i){this.type=2,this._$AH=Ue,this._$AN=void 0,this._$AA=e,this._$AB=t,this._$AM=r,this.options=i,this._$Cv=i?.isConnected??!0}get parentNode(){let e=this._$AA.parentNode;const t=this._$AM;return void 0!==t&&11===e?.nodeType&&(e=t.parentNode),e}get startNode(){return this._$AA}get endNode(){return this._$AB}_$AI(e,t=this){e=Fe(this,e,t),Pe(e)?e===Ue||null==e||""===e?(this._$AH!==Ue&&this._$AR(),this._$AH=Ue):e!==this._$AH&&e!==Ne&&this._(e):void 0!==e._$litType$?this.$(e):void 0!==e.nodeType?this.T(e):(e=>Re(e)||"function"==typeof e?.[Symbol.iterator])(e)?this.k(e):this._(e)}O(e){return this._$AA.parentNode.insertBefore(e,this._$AB)}T(e){this._$AH!==e&&(this._$AR(),this._$AH=this.O(e))}_(e){this._$AH!==Ue&&Pe(this._$AH)?this._$AA.nextSibling.data=e:this.T(Ce.createTextNode(e)),this._$AH=e}$(e){const{values:t,_$litType$:r}=e,i="number"==typeof r?this._$AC(e):(void 0===r.el&&(r.el=Ke.createElement(Ve(r.h,r.h[0]),this.options)),r);if(this._$AH?._$AD===i)this._$AH.p(t);else{const e=new We(i,this),r=e.u(this.options);e.p(t),this.T(r),this._$AH=e}}_$AC(e){let t=De.get(e.strings);return void 0===t&&De.set(e.strings,t=new Ke(e)),t}k(e){Re(this._$AH)||(this._$AH=[],this._$AR());const t=this._$AH;let r,i=0;for(const o of e)i===t.length?t.push(r=new Ye(this.O(Ee()),this.O(Ee()),this,this.options)):r=t[i],r._$AI(o),i++;i<t.length&&(this._$AR(r&&r._$AB.nextSibling,i),t.length=i)}_$AR(e=this._$AA.nextSibling,t){for(this._$AP?.(!1,!0,t);e&&e!==this._$AB;){const t=e.nextSibling;e.remove(),e=t}}setConnected(e){void 0===this._$AM&&(this._$Cv=e,this._$AP?.(e))}}class Ge{get tagName(){return this.element.tagName}get _$AU(){return this._$AM._$AU}constructor(e,t,r,i,o){this.type=1,this._$AH=Ue,this._$AN=void 0,this.element=e,this.name=t,this._$AM=i,this.options=o,r.length>2||""!==r[0]||""!==r[1]?(this._$AH=Array(r.length-1).fill(new String),this.strings=r):this._$AH=Ue}_$AI(e,t=this,r,i){const o=this.strings;let n=!1;if(void 0===o)e=Fe(this,e,t,0),n=!Pe(e)||e!==this._$AH&&e!==Ne,n&&(this._$AH=e);else{const i=e;let a,s;for(e=o[0],a=0;a<o.length-1;a++)s=Fe(this,i[r+a],t,a),s===Ne&&(s=this._$AH[a]),n||=!Pe(s)||s!==this._$AH[a],s===Ue?e=Ue:e!==Ue&&(e+=(s??"")+o[a+1]),this._$AH[a]=s}n&&!i&&this.j(e)}j(e){e===Ue?this.element.removeAttribute(this.name):this.element.setAttribute(this.name,e??"")}}class qe extends Ge{constructor(){super(...arguments),this.type=3}j(e){this.element[this.name]=e===Ue?void 0:e}}class Je extends Ge{constructor(){super(...arguments),this.type=4}j(e){this.element.toggleAttribute(this.name,!!e&&e!==Ue)}}class Qe extends Ge{constructor(e,t,r,i,o){super(e,t,r,i,o),this.type=5}_$AI(e,t=this){if((e=Fe(this,e,t,0)??Ue)===Ne)return;const r=this._$AH,i=e===Ue&&r!==Ue||e.capture!==r.capture||e.once!==r.once||e.passive!==r.passive,o=e!==Ue&&(r===Ue||i);i&&this.element.removeEventListener(this.name,this,r),o&&this.element.addEventListener(this.name,this,e),this._$AH=e}handleEvent(e){"function"==typeof this._$AH?this._$AH.call(this.options?.host??this.element,e):this._$AH.handleEvent(e)}}class Ze{constructor(e,t,r){this.element=e,this.type=6,this._$AN=void 0,this._$AM=t,this.options=r}get _$AU(){return this._$AM._$AU}_$AI(e){Fe(this,e)}}const Xe={I:Ye},et=be.litHtmlPolyfillSupport;et?.(Ke,Ye),(be.litHtmlVersions??=[]).push("3.2.1");
/**
 * @license
 * Copyright 2017 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */
let tt=class extends ye{constructor(){super(...arguments),this.renderOptions={host:this},this._$Do=void 0}createRenderRoot(){const e=super.createRenderRoot();return this.renderOptions.renderBefore??=e.firstChild,e}update(e){const t=this.render();this.hasUpdated||(this.renderOptions.isConnected=this.isConnected),super.update(e),this._$Do=((e,t,r)=>{const i=r?.renderBefore??t;let o=i._$litPart$;if(void 0===o){const e=r?.renderBefore??null;i._$litPart$=o=new Ye(t.insertBefore(Ee(),e),e,void 0,r??{})}return o._$AI(e),o})(t,this.renderRoot,this.renderOptions)}connectedCallback(){super.connectedCallback(),this._$Do?.setConnected(!0)}disconnectedCallback(){super.disconnectedCallback(),this._$Do?.setConnected(!1)}render(){return Ne}};tt._$litElement$=!0,tt.finalized=!0,globalThis.litElementHydrateSupport?.({LitElement:tt});const rt=globalThis.litElementPolyfillSupport;rt?.({LitElement:tt}),(globalThis.litElementVersions??=[]).push("4.1.1");
/**
 * @license
 * Copyright 2017 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */
const it=e=>(t,r)=>{void 0!==r?r.addInitializer((()=>{customElements.define(e,t)})):customElements.define(e,t)}
/**
 * @license
 * Copyright 2017 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */,ot={attribute:!0,type:String,converter:fe,reflect:!1,hasChanged:ve},nt=(e=ot,t,r)=>{const{kind:i,metadata:o}=r;let n=globalThis.litPropertyMetadata.get(o);if(void 0===n&&globalThis.litPropertyMetadata.set(o,n=new Map),n.set(r.name,e),"accessor"===i){const{name:i}=r;return{set(r){const o=t.get.call(this);t.set.call(this,r),this.requestUpdate(i,o,e)},init(t){return void 0!==t&&this.P(i,void 0,e),t}}}if("setter"===i){const{name:i}=r;return function(r){const o=this[i];t.call(this,r),this.requestUpdate(i,o,e)}}throw Error("Unsupported decorator location: "+i)};function at(e){return(t,r)=>"object"==typeof r?nt(e,t,r):((e,t,r)=>{const i=t.hasOwnProperty(r);return t.constructor.createProperty(r,i?{...e,wrapped:!0}:e),i?Object.getOwnPropertyDescriptor(t,r):void 0})(e,t,r)
/**
 * @license
 * Copyright 2017 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */}function st(e){return at({...e,state:!0,attribute:!1})}var lt=re`
  * {
    padding: 0;
    margin: 0;
    box-sizing: border-box;
  }

  .skeletonBase {
    transition: opacity 200ms;
    opacity: 0;
  }
  .skeletonBase.skeleton {
    opacity: 1;
    position: relative;
    /* other styles */
    width: fit-content;
    transition: opacity 200ms;
  }

  .skeletonCache {
    font-style: italic;
    opacity: 0.85;
  }

  .skeletonBase.skeletonLoaded:not(.skeleton) {
    opacity: 1;
    transition: opacity 0ms !important;
  }
  .skeletonBase.skeleton::before {
    content: "";
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: linear-gradient(45deg, var(--primary-100), var(--primary-200));
    transition: opacity 200ms;
    background-size: 200% 200%; /* Makes the animation smoother */

    pointer-events: none;
    animation: gradient-animation 500ms infinite;
    z-index: 1;
    border-radius: 0.5rem;
    opacity: 1;
  }

  .skeletonBase.skeleton.skeletonLoaded::before {
    opacity: 0;
  }

  .skeletonBase.skeleton:not(.skeletonLoaded) > *,
  .skeletonBase.skeletonResolved.skeleton:not(.skeletonLoaded) {
    color: transparent !important;
  }

  .skeletonPointerOnly {
    pointer-events: none; /* Disables all hover and click events */
  }

  @keyframes gradient-animation {
    0% {
      background-position: 0% 50%;
    }
    50% {
      background-position: 100% 50%;
    }
    100% {
      background-position: 0% 50%;
    }
  }
  header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding-bottom: 1rem;
    flex-wrap: wrap;
    gap: 0rem 1rem;
  }

  header > div {
    display: flex;
    gap: 1rem;
    flex-direction: column;
  }

  header > div.float {
    flex-direction: row;
    align-items: center;
    flex-wrap: wrap;
  }

  .btn {
    text-align: center;
    border: 0;
    padding: 0.75rem 1rem;
    text-decoration: none;
    border-radius: 0.5rem;
    transition: 200ms;
    cursor: pointer;
    font-size: 1rem;
  }
  .btn:not(.sm) {
    font-weight: 600;
  }
  .btn:not(.tertiary):not(:disabled) {
    -webkit-box-shadow: 0px 1px 5px 0px rgba(0, 0, 0, 0.15);
    -moz-box-shadow: 0px 1px 5px 0px rgba(0, 0, 0, 0.15);
    box-shadow: 0px 1px 5px 0px rgba(0, 0, 0, 0.15);
  }
  .btn.big {
    padding: 1rem 1rem;
    font-size: 1rem;
  }
  .btn:disabled {
    cursor: not-allowed;
  }

  .btn.primary {
    background-color: var(--primary-500);
    color: #ffffff;
  }
  .btn.primary.danger {
    background-color: var(--danger-600);
  }
  .btn.primary:hover {
    background-color: var(--primary-600);
  }
  .btn.primary.danger:hover {
    background-color: var(--danger-700);
  }
  .btn.primary:active {
    background-color: var(--primary-900);
  }
  .btn.primary.danger:active {
    background-color: var(--danger-800);
  }
  .btn.primary:disabled {
    background-color: var(--primary-100);
    color: var(--primary-600);
  }
  .btn.primary.danger:disabled {
    background-color: var(--danger-100);
    color: var(--danger-600);
  }
  .btn.rounded {
    border-radius: 2rem;
    padding: 0.55rem 2rem;
    font-weight: normal;
  }

  .btn.tertiary {
    padding: 0.25rem 0.5rem;
    border-bottom: 1px solid var(--primary-500);
    color: var(--primary-500);
    border-radius: 0;
    background: 0;
  }
  .btn.tertiary.danger {
    border-bottom: 1px solid var(--danger-500);
    color: var(--danger-500);
  }
  .btn.tertiary:hover {
    border-bottom: 1px solid var(--primary-600);
    color: var(--primary-600);
  }
  .btn.tertiary.danger:hover {
    border-bottom: 1px solid var(--danger-600);
    color: var(--danger-600);
  }
  .btn.tertiary:active {
    border-bottom: 1px solid var(--primary-700);
    color: var(--primary-700);
  }
  .btn.tertiary.danger:active {
    border-bottom: 1px solid var(--danger-700);
    color: var(--danger-700);
  }
  .btn.tertiary:disabled {
    border-bottom: 1px solid var(--primary-200);
    color: var(--primary-200);
  }
  .btn.tertiary.danger:disabled {
    border-bottom: 1px solid var(--danger-200);
    color: var(--danger-200);
  }

  .btn.secondary {
    outline: 1px solid var(--primary-300);
    background: #fff;
    color: var(--primary-600);
  }
  .btn.secondary.danger {
    outline: 1px solid var(--danger-300);
    color: var(--danger-500);
  }
  .btn.secondary:not(:disabled):hover {
    outline: 1px solid var(--primary-200);
    background-color: var(--primary-50);
  }
  .btn.secondary.danger:not(:disabled):hover {
    outline: 1px solid var(--danger-200);
    background-color: var(--danger-50);
  }
  .btn.secondary:not(:disabled):active {
    outline: 1px solid var(--primary-100);
    background-color: var(--primary-200);
  }
  .btn.secondary.danger:not(:disabled):active {
    outline: 1px solid var(--danger-100);
    background-color: var(--danger-200);
  }
  .btn.secondary:disabled {
    outline: 1px solid var(--primary-200);
    color: var(--primary-300);
  }
  .btn.secondary.danger:disabled {
    outline: 1px solid var(--danger-200);
    color: var(--danger-300);
  }

  .btn:not(:active):focus-visible {
    outline: 3px solid var(--primary-300);
  }

  .btn.danger:not(:active):focus-visible {
    outline: 3px solid var(--danger-900);
  }

  header.x {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }
  header.x button {
    margin-bottom: 0.05rem;
  }
  div#header h1,
  header h1,
  .sectionHeader {
    font-size: 1.35rem;
    font-weight: 500;
    color: var(--primary-800);
    margin-bottom: 0.75rem;
  }

  .sectionSubheader {
    font-size: 1.25rem;
    font-weight: 500;
    color: var(--primary-800);
  }

  div#header h1.split,
  header h1.split {
    display: flex;
    gap: 0.25rem;
    color: var(--primary-800);
  }

  div#header h1.split > span,
  header h1.split > span {
    color: var(--primary-950);
  }

  header + #content {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }

  .subsectionHeader {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
    margin-bottom: 1rem;
  }

  .subsectionHeader p#title {
    font-weight: 600;
    color: var(--primary-900);
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .subsectionHeader p:not(#title) {
    color: var(--gray-600);
  }

  hr {
    border: 0.5px solid var(--gray-200);
  }

  .widget {
    outline: 1px solid var(--primary-100);
    padding: 1rem;
    border-radius: 0.5rem;
    display: flex;
    gap: 0.15rem;
    flex-direction: column;
    background-color: #fff;
  }
  button.widget {
    border: 0;
    font-size: inherit;
    text-align: left;
  }
  button.widget:hover,
  button.widget:focus-visible {
    outline: 2px solid var(--primary-500);
  }
  .widget.red {
    outline: 1px solid var(--danger-400);
    background-color: var(--danger-50);
    color: var(--danger-950);
  }
  .widget.red #description {
    color: var(--danger-800);
  }
  .widget.primary {
    outline: 1px solid var(--primary-400);
    background-color: var(--primary-50);
    color: var(--danger-950);
  }
  .widget.center span.icon {
    font-size: 2rem;
    color: var(--primary-600);
    height: auto;
    width: auto;
  }

  .widget.x {
    gap: 1rem;
  }

  .widget p#title {
    font-weight: 600;
  }

  .widget.callout p#title {
    font-size: 1.5rem;
    color: var(--primary-600);
    margin-bottom: 0.25rem;
  }
  .widget.callout p#title .secondary {
    font-size: 1rem;
    color: var(--gray-600);
  }

  .widget.callout p a.secondary {
    text-decoration: none;
    color: var(--primary-600);
  }

  .widget.callout p a.secondary:hover {
    color: var(--primary-600) !important;
  }

  .widget p#description {
    font-size: 0.85rem;
    color: var(--gray-700);
  }
  .widget.callout p#description {
    font-size: 0.95rem;
    color: var(--gray-800);
  }
  .widget.callout p#description .secondary {
    font-size: 0.85rem;
    color: var(--gray-600);
  }

  a[target="_blank"] {
    color: var(--primary-600);
  }

  a[target="_blank"]:hover {
    color: var(--primary-500);
  }

  a[target="_blank"]::after {
    display: inline-block;
    content: " ";
    background-repeat: no-repeat;
    width: 1ch;
    height: 1ch;
    background-size: contain;
    margin-left: 0.5rem;
  }

  .widget.center {
    align-items: center;
    justify-content: center;
  }

  button.widget.clickable {
    transition: 200ms;
    cursor: pointer;
    border: 0;
    background: 0;
    background-color: #fff;
  }
  button.widget.clickable:hover,
  button.widget.clickable:focus-visible {
    background-color: var(--gray-50);
  }
  button.widget.clickable:active {
    background-color: var(--gray-100);
    outline: 1px solid var(--gray-200);
  }
  button.widget.clickable:focus-visible {
    outline: 2px solid var(--primary-600);
  }
  button.widget.clickable p#title {
    font-size: 1rem;
    font-weight: 600;
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
    color: var(--gray-600);
    font-size: 0.85rem;
    margin: 0.25rem 0 0.5rem 0;
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
    border: 1px solid var(--gray-300);
    width: 100%;
    font-size: 1rem;
  }

  .input-container input:focus {
    outline: 2px solid var(--primary-500);
  }

  .input-container p#error {
    margin-top: 0.25rem;
    color: var(--danger-600);
    font-size: 0.85rem;
    display: none;
  }

  .input-container input:invalid + p#error {
    display: inherit;
  }

  .input-container input:invalid {
    outline: 2px solid var(--danger-600);
  }

  .list .item {
    background: 0;
    width: 100%;
    border: 0;
    text-align: left;
  }

  .list button.item,
  .list a.item {
    text-decoration: none;
    color: #000;
    cursor: pointer;
    transition: 200ms;
    margin: 0;
    font-size: 1rem;
  }

  .list button.item:hover,
  .list button.item:focus-visible,
  .list a.item:hover,
  .list a.item:focus-visible,
  .list button.item.holdhover {
    background-color: var(--primary-50);
    padding-left: 1rem;
    padding-right: 1rem;
    border-radius: 0.25rem;
  }

  .list button.item:active,
  .list a.item:active {
    background-color: var(--primary-100);
  }

  .list button.item:focus-visible,
  .list a.item:focus-visible {
    outline: 2px solid var(--primary-600);
  }

  .list .item {
    padding: 1rem 0rem;
  }

  .list .item p#title {
    font-size: 1rem;
    margin-bottom: 0.25rem;
    font-weight: 600;
    color: var(--primary-900);
  }

  .list .item p#description {
    font-size: 0.85rem;
    color: var(--gray-600);
  }

  .list .item.row {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
  }

  .list .item:not(:first-of-type) {
    border-top: 1px solid var(--gray-200);
  }

  .list .section .item {
    border-left: 4px solid var(--itemColor, var(--primary-400));
    border-top: 0;
    padding: 0.5rem 0.5rem;
  }

  .list .section .item:not(:last-of-type) {
    box-shadow: inset 0 -1px 0 0 var(--gray-200); /* fakes border-bottom */
  }

  .list .section .header {
    font-weight: 600;
    font-size: 1rem;
    color: var(--primary-900);
    margin-bottom: 0.5rem;
  }

  .list .section:not(:first-of-type) .header {
    margin-top: 1rem;
    padding-top: 1rem;
    border-top: 1px solid var(--gray-200);
  }

  .grid {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .grid > div {
    display: flex;
    flex-direction: column;
  }

  .grid > div > .widget {
    flex-grow: 1;
  }

  .flex-row {
    display: flex;
    gap: 0.75rem;
    align-items: center;
  }

  .jumper {
    left: -100000dvw;
    position: absolute;
    transition: 0ms !important;
  }

  .jumper:focus {
    position: relative;
    left: 0;
    transition: 0ms !important;
  }

  @media (min-width: 680px) {
    .subsectionHeader {
      flex-direction: row;
      gap: 0.75rem;
    }
    .grid {
      display: grid;
      grid-gap: 1rem;
    }

    .grid.two {
      grid-template-columns: 1fr 1fr;
    }

    .grid.three {
      grid-template-columns: 1fr 1fr 1fr;
    }

    .grid.four {
      grid-template-columns: 1fr 1fr;
    }

    .widget.x {
      flex-direction: row;
      gap: 1rem;
      justify-content: space-between;
      align-items: center;
    }
    .widget.x > div {
      display: flex;
      flex-direction: column;
      gap: 0.15rem;
    }
  }

  @media (min-width: 1200px) {
    .grid.four {
      grid-template-columns: 1fr 1fr 1fr 1fr;
    }
  }
`
/**
 * @license
 * Copyright 2020 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */;const{I:ct}=Xe,dt=()=>document.createComment(""),pt=(e,t,r)=>{const i=e._$AA.parentNode,o=void 0===t?e._$AB:t._$AA;if(void 0===r){const t=i.insertBefore(dt(),o),n=i.insertBefore(dt(),o);r=new ct(t,n,e,e.options)}else{const t=r._$AB.nextSibling,n=r._$AM,a=n!==e;if(a){let t;r._$AQ?.(e),r._$AM=e,void 0!==r._$AP&&(t=e._$AU)!==n._$AU&&r._$AP(t)}if(t!==o||a){let e=r._$AA;for(;e!==t;){const t=e.nextSibling;i.insertBefore(e,o),e=t}}}return r},ht=(e,t,r=e)=>(e._$AI(t,r),e),ut={},mt=e=>{e._$AP?.(!1,!0);let t=e._$AA;const r=e._$AB.nextSibling;for(;t!==r;){const e=t.nextSibling;t.remove(),t=e}},ft=1,vt=2,gt=e=>(...t)=>({_$litDirective$:e,values:t});
/**
 * @license
 * Copyright 2017 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */class yt{constructor(e){}get _$AU(){return this._$AM._$AU}_$AT(e,t,r){this._$Ct=e,this._$AM=t,this._$Ci=r}_$AS(e,t){return this.update(e,t)}update(e,t){return this.render(...t)}}
/**
 * @license
 * Copyright 2017 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */const bt=(e,t)=>{const r=e._$AN;if(void 0===r)return!1;for(const e of r)e._$AO?.(t,!1),bt(e,t);return!0},wt=e=>{let t,r;do{if(void 0===(t=e._$AM))break;r=t._$AN,r.delete(e),e=t}while(0===r?.size)},xt=e=>{for(let t;t=e._$AM;e=t){let r=t._$AN;if(void 0===r)t._$AN=r=new Set;else if(r.has(e))break;r.add(e),_t(t)}};function $t(e){void 0!==this._$AN?(wt(this),this._$AM=e,xt(this)):this._$AM=e}function kt(e,t=!1,r=0){const i=this._$AH,o=this._$AN;if(void 0!==o&&0!==o.size)if(t)if(Array.isArray(i))for(let e=r;e<i.length;e++)bt(i[e],!1),wt(i[e]);else null!=i&&(bt(i,!1),wt(i));else bt(this,e)}const _t=e=>{e.type==vt&&(e._$AP??=kt,e._$AQ??=$t)};class At extends yt{constructor(){super(...arguments),this._$AN=void 0}_$AT(e,t,r){super._$AT(e,t,r),xt(this),this.isConnected=e._$AU}_$AO(e,t=!0){e!==this.isConnected&&(this.isConnected=e,e?this.reconnected?.():this.disconnected?.()),t&&(bt(this,e),wt(this))}setValue(e){if((e=>void 0===e.strings)(this._$Ct))this._$Ct._$AI(e,this);else{const t=[...this._$Ct._$AH];t[this._$Ci]=e,this._$Ct._$AI(t,this,0)}}disconnected(){}reconnected(){}}
/**
 * @license
 * Copyright 2020 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */const Ct=()=>new Et;class Et{}const Pt=new WeakMap,Rt=gt(class extends At{render(e){return Ue}update(e,[t]){const r=t!==this.Y;return r&&void 0!==this.Y&&this.rt(void 0),(r||this.lt!==this.ct)&&(this.Y=t,this.ht=e.options?.host,this.rt(this.ct=e.element)),Ue}rt(e){if(this.isConnected||(e=void 0),"function"==typeof this.Y){const t=this.ht??globalThis;let r=Pt.get(t);void 0===r&&(r=new WeakMap,Pt.set(t,r)),void 0!==r.get(this.Y)&&this.Y.call(this.ht,void 0),r.set(this.Y,e),void 0!==e&&this.Y.call(this.ht,e)}else this.Y.value=e}get lt(){return"function"==typeof this.Y?Pt.get(this.ht??globalThis)?.get(this.Y):this.Y?.value}disconnected(){this.lt===this.ct&&this.rt(void 0)}reconnected(){this.rt(this.ct)}}),Ot=(e,t,r)=>{const i=new Map;for(let o=t;o<=r;o++)i.set(e[o],o);return i},St=gt(class extends yt{constructor(e){if(super(e),e.type!==vt)throw Error("repeat() can only be used in text expressions")}dt(e,t,r){let i;void 0===r?r=t:void 0!==t&&(i=t);const o=[],n=[];let a=0;for(const t of e)o[a]=i?i(t,a):a,n[a]=r(t,a),a++;return{values:n,keys:o}}render(e,t,r){return this.dt(e,t,r).values}update(e,[t,r,i]){const o=(e=>e._$AH)(e),{values:n,keys:a}=this.dt(t,r,i);if(!Array.isArray(o))return this.ut=a,n;const s=this.ut??=[],l=[];let c,d,p=0,h=o.length-1,u=0,m=n.length-1;for(;p<=h&&u<=m;)if(null===o[p])p++;else if(null===o[h])h--;else if(s[p]===a[u])l[u]=ht(o[p],n[u]),p++,u++;else if(s[h]===a[m])l[m]=ht(o[h],n[m]),h--,m--;else if(s[p]===a[m])l[m]=ht(o[p],n[m]),pt(e,l[m+1],o[p]),p++,m--;else if(s[h]===a[u])l[u]=ht(o[h],n[u]),pt(e,o[p],o[h]),h--,u++;else if(void 0===c&&(c=Ot(a,u,m),d=Ot(s,p,h)),c.has(s[p]))if(c.has(s[h])){const t=d.get(a[u]),r=void 0!==t?o[t]:null;if(null===r){const t=pt(e,o[p]);ht(t,n[u]),l[u]=t}else l[u]=ht(r,n[u]),pt(e,o[p],r),o[t]=null;u++}else mt(o[h]),h--;else mt(o[p]),p++;for(;u<=m;){const t=pt(e,l[m+1]);ht(t,n[u]),l[u++]=t}for(;p<=h;){const e=o[p++];null!==e&&mt(e)}return this.ut=a,((e,t=ut)=>{e._$AH=t})(e,l),Ne}}),zt=gt(class extends yt{constructor(e){if(super(e),e.type!==ft||"class"!==e.name||e.strings?.length>2)throw Error("`classMap()` can only be used in the `class` attribute and must be the only part in the attribute.")}render(e){return" "+Object.keys(e).filter((t=>e[t])).join(" ")+" "}update(e,[t]){if(void 0===this.st){this.st=new Set,void 0!==e.strings&&(this.nt=new Set(e.strings.join(" ").split(/\s/).filter((e=>""!==e))));for(const e in t)t[e]&&!this.nt?.has(e)&&this.st.add(e);return this.render(t)}const r=e.element.classList;for(const e of this.st)e in t||(r.remove(e),this.st.delete(e));for(const e in t){const i=!!t[e];i===this.st.has(e)||this.nt?.has(e)||(i?(r.add(e),this.st.add(e)):(r.remove(e),this.st.delete(e)))}return Ne}});
/**
 * @license
 * Copyright 2018 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */var Tt=function(e,t,r,i){var o,n=arguments.length,a=n<3?t:null===i?i=Object.getOwnPropertyDescriptor(t,r):i;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(e,t,r,i);else for(var s=e.length-1;s>=0;s--)(o=e[s])&&(a=(n<3?o(a):n>3?o(t,r,a):o(t,r))||a);return n>3&&a&&Object.defineProperty(t,r,a),a};let Mt=class extends tt{constructor(){super(...arguments),this.width=20,this.enabled=!1,this.label="Click to check",this.userSelectable=!0}handleClick(){this.enabled=!this.enabled;const e=new CustomEvent("checked",{detail:{status:this.enabled}});this.dispatchEvent(e)}get value(){return this.enabled?"true":"false"}set value(e){this.enabled="true"===e}render(){return Ie`<button
      @click=${this.handleClick}
      aria-label="${this.label}"
      style="height: ${this.width}px;"
      ?disabled=${!this.userSelectable}
    >
      <svg
        width="${this.width}"
        height="${this.width}"
        class=${zt({active:this.enabled})}
      >
        <circle
          id="outer"
          cx="${this.width/2}"
          cy="${this.width/2}"
          r="${this.width/2-1}"
          stroke-width="2"
        />
        <circle
          id="inner"
          cx="${this.width/2}"
          cy="${this.width/2}"
          r="${this.width/4-1}"
        />
      </svg>
    </button>`}};Mt.styles=[re`
      :host {
        --using-border: var(--fl-checkbox-border, var(--gray-400, #989898));
        --using-background: var(
          --fl-checkbox-unchecked-background,
          var(--gray-50, #f8f8f8)
        );
        --inner: transparent;
        height: 20px;
      }

      button {
        height: var(--height, 20px);
        margin: 0;
        padding: 0;
        box-sizing: border-box;
        width: fit-content;
        background: 0;
        border: 0;
        border-radius: 100%;
      }

      svg {
        transition: 250ms;
      }

      :host(.emulate-hover) button svg,
      button:hover svg,
      button:focus-visible svg {
        transition: 250ms;
        --using-border: var(
          --fl-checkbox-border-hover,
          var(--primary-400, #59a2ff)
        );
        --inner: var(--using-border);
        cursor: pointer;
      }

      svg.active circle#outer {
        --using-border: var(
          --fl-checkbox-border-active,
          var(--primary-600, #1a5cf4)
        );
        --using-background: #fff;
      }

      svg.active circle#inner {
        --inner: var(--fl-checkbox-border-active, var(--primary-600, #1a5cf4));
      }

      svg circle#outer {
        fill: var(--using-background);
        stroke: var(--using-border);
        transition:
          fill 250ms,
          stroke 250ms;
      }

      svg circle#inner {
        fill: var(--inner);
        transition: fill 250ms;
      }

      svg circle#outer {
        fill: var(--using-background);
        stroke: var(--using-border);
      }
      svg circle#inner {
        fill: var(--inner);
      }
    `],Tt([at()],Mt.prototype,"width",void 0),Tt([at()],Mt.prototype,"enabled",void 0),Tt([at()],Mt.prototype,"label",void 0),Tt([at()],Mt.prototype,"userSelectable",void 0),Mt=Tt([it("radio-input")],Mt);var Lt=function(e,t,r,i){var o,n=arguments.length,a=n<3?t:null===i?i=Object.getOwnPropertyDescriptor(t,r):i;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(e,t,r,i);else for(var s=e.length-1;s>=0;s--)(o=e[s])&&(a=(n<3?o(a):n>3?o(t,r,a):o(t,r))||a);return n>3&&a&&Object.defineProperty(t,r,a),a};let jt=class extends tt{constructor(){super(...arguments),this.default="",this.activeChoice="",this.hoveringChoiceKey="",this.choices=[{key:"small",title:"Small",description:"Small T-Shirt Size"},{key:"medium",title:"Medium",description:"Medium T-Shirt Size"},{key:"large",title:"Large",description:"A big big t shirt"}]}get value(){return this.activeChoice}set value(e){this.activeChoice=e}handleClick(e){this.activeChoice=e.key;const t=new CustomEvent("selected",{detail:{key:e.key}});this.dispatchEvent(t)}render(){return Ie`<div id="container">
      ${St(this.choices,(e=>e.key),(e=>Ie`<div
            class="choice"
            @click=${()=>{this.handleClick(e)}}
            @mouseenter=${()=>{this.hoveringChoiceKey=e.key}}
            @mouseleave=${()=>{this.hoveringChoiceKey=""}}
          >
            <div>
              <radio-input
                class="${zt({"emulate-hover":this.hoveringChoiceKey===e.key})}"
                .label=${`${e.key}: ${e.description}`}
                .enabled=${this.activeChoice===e.key||""===this.activeChoice&&this.default===e.key}
              ></radio-input>
              <p id="title">${e.title}</p>
            </div>
            ${e.description?Ie`<p id="description">${e.description}</p>`:""}
          </div>`))}
    </div>`}};jt.styles=[re`
      #container {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
      }
      .choice {
        cursor: pointer;
      }
      .choice div {
        display: flex;
        gap: 0.5rem;
        align-items: center;
      }

      p {
        margin: 0;
        padding: 0;
      }

      .choice p#title {
        font-weight: 600;
        font-size: 1rem;
      }
      .choice p#description {
        margin-top: 0.1rem;
        font-size: 0.85rem;
        margin-left: calc(20px + 0.5rem);
      }
    `],Lt([at()],jt.prototype,"default",void 0),Lt([st()],jt.prototype,"activeChoice",void 0),Lt([st()],jt.prototype,"hoveringChoiceKey",void 0),Lt([at()],jt.prototype,"choices",void 0),jt=Lt([it("radio-selector")],jt);var Bt=function(e,t,r,i){var o,n=arguments.length,a=n<3?t:null===i?i=Object.getOwnPropertyDescriptor(t,r):i;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(e,t,r,i);else for(var s=e.length-1;s>=0;s--)(o=e[s])&&(a=(n<3?o(a):n>3?o(t,r,a):o(t,r))||a);return n>3&&a&&Object.defineProperty(t,r,a),a};let It=class extends tt{constructor(){super(...arguments),this.open=!1,this.wantedWidth="24rem",this.modaltitle="Needs Title",this.lastFocus=null,this._handleKeyDown=e=>{"Escape"===e.key&&this.open&&this.close()}}close(){this.dispatchEvent(new CustomEvent("close")),this.lastFocus&&this.lastFocus.focus()}connectedCallback(){super.connectedCallback(),window.addEventListener("keydown",this._handleKeyDown)}disconnectedCallback(){window.removeEventListener("keydown",this._handleKeyDown),super.disconnectedCallback()}updated(e){var t;if(super.updated(e),e.has("open")&&this.open&&e.get("open")!==this.open){const e=null===(t=this.shadowRoot)||void 0===t?void 0:t.getElementById("content");e&&(this.lastFocus=document.activeElement,e.focus())}e.has("open")&&!this.open&&e.get("open")!==this.open&&this.lastFocus&&this.lastFocus.focus()}render(){return this.open?Ie`
      <div>
        <div id="backdrop" @click=${this.close}></div>
        <div
          id="content"
          tabindex="-1"
          style="--wantedWidth: ${this.wantedWidth}"
        >
          <div id="header">
            <p>${this.modaltitle}</p>

            <button
              @click=${this.close}
              aria-label="close modal"
              title="Close modal"
            >
              &times;
            </button>
          </div>

          <slot><p>Please fill slot.</p></slot>
        </div>
      </div>
    `:Ie``}};It.styles=re`
    * {
      box-sizing: border-box;
      padding: 0;
      margin: 0;
    }
    :host {
      position: absolute;
    }
    :host > div {
      position: fixed;
      top: 0;
      left: 0;

      --backdrop: rgba(17, 17, 17, 0.5);

      --width: 100vw;
      --width: 100dvw;

      --height: 100vh;
      --height: 100dvh;

      width: var(--width);
      height: var(--height);

      display: flex;
      align-items: center;
      justify-content: center;

      z-index: 10000;
    }

    #backdrop {
      position: fixed;
      top: 0;
      left: 0;

      width: var(--width);
      height: var(--height);

      background-color: var(--backdrop);
    }

    #content {
      --wantedWidth: 4rem;
      position: fixed;
      z-index: 10001;
      flex-grow: 1;
      top: 0;
      left: 0;
      width: var(--width);
      height: var(--height);
      background-color: #fefefe;
      padding: 2rem;
      border-radius: 0.5rem;
      height: 100%;
      display: flex;
      flex-direction: column;
    }

    #content #header {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      gap: 1rem;
      margin-bottom: 1rem;
    }

    #content #header p {
      font-weight: 600;
      font-size: 1.25rem;
    }

    #content #header button {
      background: unset;
      border: 0;
      font-size: 1.15rem;
      cursor: pointer;
      color: var(--gray-600, #656565);
    }

    @media (min-width: 650px) {
      #content {
        position: relative;
        max-width: var(--wantedWidth);
        height: fit-content;
      }
    }
  `,Bt([at({type:Boolean})],It.prototype,"open",void 0),Bt([at({type:String})],It.prototype,"wantedWidth",void 0),Bt([at()],It.prototype,"modaltitle",void 0),It=Bt([it("modal-component")],It);var Nt=function(e,t,r,i){var o,n=arguments.length,a=n<3?t:null===i?i=Object.getOwnPropertyDescriptor(t,r):i;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(e,t,r,i);else for(var s=e.length-1;s>=0;s--)(o=e[s])&&(a=(n<3?o(a):n>3?o(t,r,a):o(t,r))||a);return n>3&&a&&Object.defineProperty(t,r,a),a};let Ut=class extends tt{constructor(){super(...arguments),this.disabled=!1,this.expectLoad=!1,this.loading=!1,this.loadingText=""}render(){return Ie`<button
      ?disabled=${this.disabled||this.loading}
      class=${zt({loading:this.loading})}
      @click=${()=>{this.dispatchEvent(new Event("fl-click"))}}
    >
      ${this.expectLoad?Ie`
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
          `:Ie``}
      ${this.loading&&""!==this.loadingText?this.loadingText:Ie`<slot></slot>`}
    </button>`}};Ut.styles=re`
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
  `,Nt([at()],Ut.prototype,"disabled",void 0),Nt([at()],Ut.prototype,"expectLoad",void 0),Nt([at()],Ut.prototype,"loading",void 0),Nt([at()],Ut.prototype,"loadingText",void 0),Ut=Nt([it("button-component")],Ut);const Dt=re`
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
    color: var(--gray-600);
    font-size: 0.85rem;
    margin: 0.25rem 0 0.5rem 0;
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
    border: 1px solid var(--gray-300);
    width: 100%;
    font-size: 1rem;
  }

  .input-container input:focus {
    outline: 2px solid var(--primary-500);
  }

  .input-container p#error {
    margin-top: 0.25rem;
    color: var(--danger-600);
    font-size: 0.85rem;
    display: none;
  }

  .input-container time-picker-component:invalid + p#error,
  .input-container :user-invalid + p#error {
    display: inherit;
  }

  .input-container time-picker-component:invalid,
  .input-container :user-invalid {
    border-radius: 0.5rem;
    outline: 2px solid var(--danger-600);
  }

  .checkbox-container {
    display: flex;
    align-items: center;
    gap: 1rem;
  }
`;var Ht=function(e,t,r,i){var o,n=arguments.length,a=n<3?t:null===i?i=Object.getOwnPropertyDescriptor(t,r):i;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(e,t,r,i);else for(var s=e.length-1;s>=0;s--)(o=e[s])&&(a=(n<3?o(a):n>3?o(t,r,a):o(t,r))||a);return n>3&&a&&Object.defineProperty(t,r,a),a};let Vt=class extends tt{constructor(){super(),this.open=!1,this.promptID="",this.promptType="text",this.promptTitle="",this.promptDescription="",this.promptValue="",this.promptPattern="",this.promptPatternError="",this.promptRadioChoices=[],this.promptRadioDefault=void 0,this.inputRef=Ct(),this.openPrompt=e=>{let{detail:t}=e;this.promptID=t.id,this.promptType=t.type,this.promptTitle=t.title,this.promptDescription=t.description,this.promptPattern=t.pattern,this.promptPatternError=t.patternError,this.promptRadioChoices=t.radioOptions,this.promptRadioDefault=t.radioDefaultKey,this.open=!0}}connectedCallback(){super.connectedCallback(),window.addEventListener("fl-prompt",this.openPrompt)}disconnectedCallback(){window.removeEventListener("fl-prompt",this.openPrompt),super.disconnectedCallback()}handleChange(e){this.promptValue=e.target.value}submitPrompt(){window.dispatchEvent(new CustomEvent(`fl-response-${this.promptID}`,{detail:{id:this.promptID,value:this.promptValue}})),this.close()}close(e){e&&window.dispatchEvent(new CustomEvent(`fl-response-${this.promptID}`,{detail:{canceled:!0}})),this.reset()}reset(){var e;this.open=!1,this.promptType="text",this.promptTitle="",this.promptDescription="",this.promptID="",this.promptValue="",this.promptPattern="",this.promptPatternError="",this.promptRadioChoices=void 0,this.promptRadioDefault=void 0,(null===(e=this.shadowRoot)||void 0===e?void 0:e.getElementById("input")).value=""}updated(e){super.updated(e),e.has("open")&&!this.open&&e.get("open")!==this.open&&this.reset()}getInputMode(){switch(this.promptType){case"number":return"numeric";case"decimal":return"decimal";case"tel":return"tel";default:return"text"}}getPromptArea(){var e,t;switch(this.promptType){case"textarea":return Ie`<textarea
          id="input"
          @input=${this.handleChange}
          placeholder="type something"
          style="flex-grow: 1;"
        ></textarea>`;case"radio":return Ie` <radio-selector
          id="input"
          style="margin: 0.6rem 0;"
          .default=${this.promptRadioDefault}
          .choices=${this.promptRadioChoices}
          @selected=${e=>{this.promptValue=e.detail.key}}
        ></radio-selector>`;default:return Ie`<input
            ${Rt(this.inputRef)}
            pattern=${(e=>e??Ue)
/**
 * @license
 * Copyright 2017 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */(this.promptPattern||void 0)}
            step="${"decimal"===this.promptType?"0.01":""}"
            type="${"decimal"===this.promptType?"number":this.promptType}"
            inputmode="${this.getInputMode()}"
            id="input"
            @input=${this.handleChange}
            placeholder="type something"
          />
          ${this.promptPatternError&&!(null===(t=null===(e=this.inputRef.value)||void 0===e?void 0:e.validity)||void 0===t?void 0:t.valid)?Ie` <p class="error">${this.promptPatternError}</p>`:""}`}}render(){var e,t,r;return Ie`<modal-component
      .open=${this.open}
      .modaltitle="${this.promptTitle}"
      @close=${this.close}
    >
      <div id="feedback">
        <div>
          ${""!==this.promptDescription?Ie` <p>${this.promptDescription}</p> `:void 0}
          ${this.getPromptArea()}
        </div>

        <button-component
          class="big"
          @fl-click=${this.submitPrompt}
          .disabled=${0===this.promptValue.length||void 0!==(null===(e=this.inputRef.value)||void 0===e?void 0:e.validity)&&!(null===(r=null===(t=this.inputRef.value)||void 0===t?void 0:t.validity)||void 0===r?void 0:r.valid)}
          >Submit</button-component
        >
      </div>
    </modal-component>`}};Vt.styles=[re`
      :host {
        --transition: 600ms;
      }

      * {
        padding: 0;
        margin: 0;
        box-sizing: border-box;
      }

      #feedback {
        display: flex;
        flex-direction: column;
        justify-content: space-between;
        height: 100%;
        gap: 1rem;
      }

      #feedback > div {
        display: flex;
        flex-direction: column;
        height: 100%;
      }

      #feedback > div > p {
        margin-bottom: 0.75rem;
      }

      p.error {
        margin-top: 1rem;
        color: var(--danger-600, #e01e47);
      }
    `,Dt],Ht([at()],Vt.prototype,"open",void 0),Ht([st()],Vt.prototype,"promptID",void 0),Ht([st()],Vt.prototype,"promptType",void 0),Ht([st()],Vt.prototype,"promptTitle",void 0),Ht([st()],Vt.prototype,"promptDescription",void 0),Ht([st()],Vt.prototype,"promptValue",void 0),Ht([st()],Vt.prototype,"promptPattern",void 0),Ht([st()],Vt.prototype,"promptPatternError",void 0),Ht([st()],Vt.prototype,"promptRadioChoices",void 0),Ht([st()],Vt.prototype,"promptRadioDefault",void 0),Vt=Ht([it("prompt-component")],Vt);var Kt=function(e,t,r,i){var o,n=arguments.length,a=n<3?t:null===i?i=Object.getOwnPropertyDescriptor(t,r):i;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(e,t,r,i);else for(var s=e.length-1;s>=0;s--)(o=e[s])&&(a=(n<3?o(a):n>3?o(t,r,a):o(t,r))||a);return n>3&&a&&Object.defineProperty(t,r,a),a};let Ft=class extends tt{constructor(){super(...arguments),this.open=!1,this.promptTitle="",this.promptDescription="",this.proceedText="Confirm",this.cancelText="Cancel",this.openAlert=e=>{let{detail:t}=e;this.promptTitle=t.title,this.promptDescription=t.description,this.proceedText=t.proceedButton||this.proceedText,this.cancelText=t.cancelButton||this.cancelText,this.open=!0}}connectedCallback(){super.connectedCallback(),window.addEventListener("fl-confirm",this.openAlert)}disconnectedCallback(){window.removeEventListener("fl-confirm",this.openAlert),super.disconnectedCallback()}close(){this.reset()}reset(){this.open=!1,this.promptTitle="",this.promptDescription="",this.proceedText="Confirm",this.cancelText="Cancel"}updated(e){super.updated(e),e.has("open")&&!this.open&&e.get("open")!==this.open&&this.reset()}render(){return Ie`<modal-component
      .open=${this.open}
      .modaltitle="${this.promptTitle}"
      @close=${this.close}
    >
      <div id="feedback">
        <div>
          <p>${this.promptDescription}</p>
        </div>

        <div id="buttons">
          <button-component class="big plain" @fl-click=${this.close}
            >${this.cancelText}</button-component
          >
          <button-component
            class="big"
            @fl-click=${this.close}
            style="flex-grow: 1; min-width: 12rem"
            >${this.proceedText}</button-component
          >
        </div>
      </div>
    </modal-component>`}};Ft.styles=[re`
      * {
        padding: 0;
        margin: 0;
        box-sizing: border-box;
      }

      #feedback {
        display: flex;
        flex-direction: column;
        justify-content: space-between;
        height: 100%;
        gap: 1rem;
      }

      #buttons {
        display: flex;
        gap: 2rem;
        width: 100%;
        flex-wrap: wrap;
        align-items: center;
      }
    `],Kt([at()],Ft.prototype,"open",void 0),Kt([st()],Ft.prototype,"promptTitle",void 0),Kt([st()],Ft.prototype,"promptDescription",void 0),Kt([st()],Ft.prototype,"proceedText",void 0),Kt([st()],Ft.prototype,"cancelText",void 0),Ft=Kt([it("confirm-component")],Ft);var Wt=function(e,t,r,i){var o,n=arguments.length,a=n<3?t:null===i?i=Object.getOwnPropertyDescriptor(t,r):i;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(e,t,r,i);else for(var s=e.length-1;s>=0;s--)(o=e[s])&&(a=(n<3?o(a):n>3?o(t,r,a):o(t,r))||a);return n>3&&a&&Object.defineProperty(t,r,a),a};let Yt=class extends tt{constructor(){super(...arguments),this.open=!1,this.promptTitle="",this.promptDescription=void 0,this.acknowledge=void 0,this.openAlert=e=>{let{detail:t}=e;this.promptTitle=t.title,this.promptDescription=t.description,this.acknowledge=t.acknowledgeText,this.open=!0}}connectedCallback(){super.connectedCallback(),window.addEventListener("fl-alert",this.openAlert)}disconnectedCallback(){window.removeEventListener("fl-alert",this.openAlert),super.disconnectedCallback()}close(){this.reset()}reset(){this.open=!1,this.promptTitle="",this.promptDescription="",this.acknowledge=void 0}updated(e){super.updated(e),e.has("open")&&!this.open&&e.get("open")!==this.open&&this.reset()}render(){return Ie`<modal-component
      .open=${this.open}
      .modaltitle="${this.promptTitle}"
      @close=${this.close}
    >
      <div id="feedback">
        <div>
          ${void 0!==this.promptDescription?Ie` <p>${this.promptDescription}</p> `:void 0}
        </div>

        <button-component class="big" @fl-click=${this.close}
          >${this.acknowledge?this.acknowledge:"OK"}</button-component
        >
      </div>
    </modal-component>`}};Yt.styles=[re`
      * {
        padding: 0;
        margin: 0;
        box-sizing: border-box;
      }

      #feedback {
        display: flex;
        flex-direction: column;
        justify-content: space-between;
        height: 100%;
        gap: 1rem;
      }
    `],Wt([at()],Yt.prototype,"open",void 0),Wt([st()],Yt.prototype,"promptTitle",void 0),Wt([st()],Yt.prototype,"promptDescription",void 0),Wt([st()],Yt.prototype,"acknowledge",void 0),Yt=Wt([it("alert-component")],Yt);const Gt={activity:Ie`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 37 37"><path fill-rule="evenodd" d="M8 2h21a6 6 0 0 1 6 6v21a6 6 0 0 1-6 6H8a6 6 0 0 1-6-6V8a6 6 0 0 1 6-6M0 8a8 8 0 0 1 8-8h21a8 8 0 0 1 8 8v21a8 8 0 0 1-8 8H8a8 8 0 0 1-8-8zm8.5 1a1.5 1.5 0 1 0 0 3h21a1.5 1.5 0 0 0 0-3zM7 18.5A1.5 1.5 0 0 1 8.5 17h18a1.5 1.5 0 0 1 0 3h-18A1.5 1.5 0 0 1 7 18.5M8.5 25a1.5 1.5 0 0 0 0 3h20a1.5 1.5 0 0 0 0-3z" class="primary" clip-rule="evenodd"/></svg>`,alert:Ie`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 100 100"><g clip-path="url(#a)"><path stroke-width="5" d="m52.165 17.75 34.641 60c.962 1.667-.24 3.75-2.165 3.75H15.359c-1.924 0-3.127-2.083-2.165-3.75l34.64-60c.963-1.667 3.369-1.667 4.331 0" class="primary-stroke"/><path d="M44.414 40.384A5 5 0 0 1 49.4 35h1.202a5 5 0 0 1 4.985 5.383l-1.114 14.475a4.486 4.486 0 0 1-8.945 0z" class="primary"/><circle cx="50" cy="68" r="5" class="primary"/></g><defs><clipPath id="a"><path fill="#fff" d="M0 0h100v100H0z"/></clipPath></defs></svg>`,check:Ie`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 17 15"><path stroke-width="2" d="m1 9.5 3.695 3.695a1 1 0 0 0 1.5-.098L15.5 1" class="primary-stroke"/></svg>`,checkmark:Ie`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 33 33"><circle cx="16.15" cy="16.15" r="16.15" class="primary"/><path stroke="#fff" stroke-width="3" d="m8.604 18.867 3.328 3.328a1 1 0 0 0 1.452-.04L24.3 9.962"/></svg>`,clock:Ie`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 100 100"><path stroke-width="6" d="M68.5 14.526A39.8 39.8 0 0 0 50 10c-22.091 0-40 17.909-40 40s17.909 40 40 40 40-17.909 40-40c0-8.127-1.336-14.688-5.5-21" class="secondary-stroke"/><path d="M87.255 18.607a5 5 0 1 0-7.071-7.071L45.536 46.184a5 5 0 1 0 7.07 7.07zM24.16 82.33a5 5 0 0 0-8.66-5l-5 8.66a5 5 0 1 0 8.66 5zm51.34 0a5 5 0 1 1 8.66-5l5 8.66a5 5 0 0 1-8.66 5z" class="primary"/></svg>`,cog:Ie`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 95 95"><path fill-rule="evenodd" d="M43 0a5 5 0 0 0-5 5v8.286c0 1.856-1.237 3.473-2.951 4.185-1.715.712-3.71.432-5.024-.88l-5.86-5.86a5 5 0 0 0-7.07 0l-6.365 6.363a5 5 0 0 0 0 7.072l5.86 5.86c1.313 1.312 1.593 3.308.88 5.023C16.76 36.763 15.143 38 13.287 38H5a5 5 0 0 0-5 5v9a5 5 0 0 0 5 5h8.286c1.856 0 3.473 1.237 4.185 2.951.712 1.715.432 3.71-.88 5.024l-5.86 5.86a5 5 0 0 0 0 7.07l6.363 6.364a5 5 0 0 0 7.072 0l5.86-5.86c1.312-1.312 3.308-1.592 5.023-.88S38 79.858 38 81.714V90a5 5 0 0 0 5 5h9a5 5 0 0 0 5-5v-8.286c0-1.856 1.237-3.473 2.951-4.185 1.715-.712 3.71-.432 5.024.88l5.86 5.86a5 5 0 0 0 7.07 0l6.365-6.363a5 5 0 0 0 0-7.071l-5.86-5.86c-1.313-1.313-1.593-3.308-.88-5.024.71-1.714 2.327-2.951 4.183-2.951H90a5 5 0 0 0 5-5v-9a5 5 0 0 0-5-5h-8.286c-1.856 0-3.473-1.237-4.185-2.951-.712-1.715-.432-3.71.88-5.024l5.86-5.86a5 5 0 0 0 0-7.07l-6.363-6.365a5 5 0 0 0-7.071 0l-5.86 5.86c-1.313 1.313-3.308 1.593-5.024.88C58.237 16.76 57 15.143 57 13.287V5a5 5 0 0 0-5-5zm4 62c8.284 0 15-6.716 15-15s-6.716-15-15-15-15 6.716-15 15 6.716 15 15 15" class="primary" clip-rule="evenodd"/></svg>`,email:Ie`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 100 100"><path d="M48.19 50.952 7 30a6 6 0 0 1 6-6h74a6 6 0 0 1 6 6L52.765 50.93a5 5 0 0 1-4.574.022" class="primary"/><path fill-rule="evenodd" d="M88 26H12a4 4 0 0 0-4 4v41a4 4 0 0 0 4 4h76a4 4 0 0 0 4-4V30a4 4 0 0 0-4-4m-76-4a8 8 0 0 0-8 8v41a8 8 0 0 0 8 8h76a8 8 0 0 0 8-8V30a8 8 0 0 0-8-8z" class="secondary" clip-rule="evenodd"/></svg>`,flag:Ie`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 9 13"><path d="M8 5V1l-1.175.294a10 10 0 0 1-5.588-.215L1 1v4l.237.08a10 10 0 0 0 5.588.214z" class="secondary"/><path d="M1 12.5V5m0 0V1l.237.08a10 10 0 0 0 5.588.214L8 1v4l-1.175.294a10 10 0 0 1-5.588-.215z" class="primary-stroke"/></svg>`,home:Ie`<svg xmlns="http://www.w3.org/2000/svg" class="icon-home" viewBox="0 0 24 24"><path d="M9 22H5a1 1 0 0 1-1-1V11l8-8 8 8v10a1 1 0 0 1-1 1h-4a1 1 0 0 1-1-1v-4a1 1 0 0 0-1-1h-2a1 1 0 0 0-1 1v4a1 1 0 0 1-1 1m3-9a2 2 0 1 0 0-4 2 2 0 0 0 0 4" class="primary"/><path d="m12.01 4.42-8.3 8.3a1 1 0 1 1-1.42-1.41l9.02-9.02a1 1 0 0 1 1.41 0l8.99 9.02a1 1 0 0 1-1.42 1.41z" class="secondary"/></svg>`,info:Ie`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"><circle cx="12" cy="12" r="11.5" stroke="#fff"/><path stroke-width="2" d="M13.5 18.5V13a1 1 0 0 0-1-1H10m3.5 6.5h-4m4 0h3" class="primary-stroke"/><circle cx="12.5" cy="7" r="2" class="primary"/></svg>`,note:Ie`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 100 100"><g stroke-width="6" clip-path="url(#a)"><path d="M58.657 3H18C9.716 3 3 9.716 3 18v64c0 8.284 6.716 15 15 15h64c8.284 0 15-6.716 15-15V34.629" class="primary-stroke"/><path d="M48.93 54.861 79.801 3.473a1 1 0 0 1 1.358-.35L92.707 9.79a1 1 0 0 1 .406 1.29l-.049.091L62.38 62.25 42.86 76.275z" class="secondary-stroke"/></g><defs><clipPath id="a"><path fill="#fff" d="M0 0h100v100H0z"/></clipPath></defs></svg>`,"person-group":Ie`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M12 13a3 3 0 0 1 3-3h4a3 3 0 0 1 3 3v3a1 1 0 0 1-1 1h-8a1 1 0 0 1-1-1 1 1 0 0 1-1 1H3a1 1 0 0 1-1-1v-3a3 3 0 0 1 3-3h4a3 3 0 0 1 3 3M7 9a3 3 0 1 1 0-6 3 3 0 0 1 0 6m10 0a3 3 0 1 1 0-6 3 3 0 0 1 0 6" class="secondary"/><path d="M12 13a3 3 0 1 1 0-6 3 3 0 0 1 0 6m-3 1h6a3 3 0 0 1 3 3v3a1 1 0 0 1-1 1H7a1 1 0 0 1-1-1v-3a3 3 0 0 1 3-3" class="primary"/></svg>`,"person-outline":Ie`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 30 34"><path fill-rule="evenodd" d="M20 7a5 5 0 1 1-10 0 5 5 0 0 1 10 0m2 0A7 7 0 1 1 8 7a7 7 0 0 1 14 0M.187 21.912C-.438 17.746 2.788 14 7 14l1.694 1.376a10 10 0 0 0 12.612 0L23 14c4.212 0 7.438 3.746 6.813 7.912L28 34H2zm22.38-4.983 1.09-.886a4.89 4.89 0 0 1 4.178 5.572L26.278 32H3.722L2.165 21.615a4.89 4.89 0 0 1 4.178-5.572l1.09.886a12 12 0 0 0 15.134 0" class="primary" clip-rule="evenodd"/></svg>`,person:Ie`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 30 34"><g clip-path="url(#a)"><path fill-rule="evenodd" d="M20 7a5 5 0 1 1-10 0 5 5 0 0 1 10 0m2 0A7 7 0 1 1 8 7a7 7 0 0 1 14 0M.187 21.912C-.438 17.746 2.788 14 7 14l1.694 1.376a10 10 0 0 0 12.612 0L23 14c4.212 0 7.438 3.746 6.813 7.912L28 34H2z" class="primary" clip-rule="evenodd"/></g><defs><clipPath id="a"><path fill="#fff" d="M0 0h30v34H0z"/></clipPath></defs></svg>`,"phone-disabled":Ie`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 100 100"><g clip-path="url(#a)"><rect width="31" height="11.499" x="37.69" y="4.483" stroke="#A4A4A4" stroke-width="4" rx="3" transform="rotate(60 37.69 4.483)"/><path stroke="#A4A4A4" stroke-width="4" d="M56.212 88a13 13 0 0 1-9.46-4.082l-.233-.255-14.483-16.209a13 13 0 0 1-2.514-4.191l-.109-.31L20.205 35.7a13 13 0 0 1 1.186-10.876l.196-.315 3.355-5.218 9.737 16.23c.21.348.345.735.4 1.136l.018.174.88 11.26a27 27 0 0 0 12.767 20.893l.719.426 6.43 3.689c.383.22.713.52.965.88l.103.158L65.434 88z"/><rect width="31" height="11.499" x="70.69" y="60.732" stroke="#A4A4A4" stroke-width="4" rx="3" transform="rotate(60 70.69 60.732)"/><circle cx="36.869" cy="60.869" r="19.869" class="primary"/><path fill="#fff" fill-rule="evenodd" d="M30.32 49.486a1 1 0 0 0-1.413 0l-3.683 3.682a1 1 0 0 0 0 1.415l5.908 5.907a1 1 0 0 1 0 1.414l-6.103 6.103a1 1 0 0 0 0 1.414l3.55 3.55a1 1 0 0 0 1.414 0l6.103-6.103a1 1 0 0 1 1.414 0l5.907 5.908a1 1 0 0 0 1.415 0l3.682-3.682a1 1 0 0 0 0-1.415l-5.908-5.907a1 1 0 0 1 0-1.414l6.103-6.103a1 1 0 0 0 0-1.415l-3.55-3.55a1 1 0 0 0-1.413 0l-6.104 6.104a1 1 0 0 1-1.414 0z" clip-rule="evenodd"/></g><defs><clipPath id="a"><path fill="#fff" d="M0 0h100v100H0z"/></clipPath></defs></svg>`,phone:Ie`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 100 100"><g clip-path="url(#a)"><rect width="35" height="15.499" x="37.422" y="-.249" class="secondary" rx="5" transform="rotate(60 37.422 -.25)"/><path d="M24.13 16.854a1 1 0 0 1 1.698.026l10.566 17.61c.399.664.637 1.412.698 2.184l.88 11.26a25 25 0 0 0 12.486 19.74l6.431 3.689a5 5 0 0 1 1.779 1.73l9.402 15.386A1 1 0 0 1 67.217 90H56.212a15 15 0 0 1-11.185-5.005L30.544 68.787a15 15 0 0 1-3.026-5.193L18.311 36.34a15 15 0 0 1 1.593-12.913z" class="primary"/><rect width="35" height="15.499" x="70.422" y="56" class="secondary" rx="5" transform="rotate(60 70.422 56)"/></g><defs><clipPath id="a"><path fill="#fff" d="M0 0h100v100H0z"/></clipPath></defs></svg>`,pin:Ie`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 100 100"><path fill="#fff" d="M0 0h100v100H0z"/><path fill-rule="evenodd" d="M34.825 12.39A5 5 0 0 1 39.787 8H60.94a5 5 0 0 1 4.971 4.465l5.939 55.142a5 5 0 0 1-4.971 5.535h-5.264q.036-.456.036-.923c0-2.683-.914-5.153-2.447-7.116A5 5 0 0 0 62.37 59.9l-2.89-26.696a5 5 0 0 0-4.971-4.462H46.4a5 5 0 0 0-4.963 4.386l-3.302 26.697A5 5 0 0 0 41.045 65a11.52 11.52 0 0 0-2.493 8.142h-5.551a5 5 0 0 1-4.963-5.61z" class="primary" clip-rule="evenodd"/><circle cx="49.868" cy="72" r="7" class="secondary"/><rect width="8" height="18" x="46" y="75" class="secondary" rx="3"/></svg>`,search:Ie`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"><path fill-rule="evenodd" d="M9.92 7.93a5.93 5.93 0 1 1 11.858 0 5.93 5.93 0 0 1-11.859 0M15.848 0a7.93 7.93 0 0 0-6.27 12.785L.293 22.07l1.414 1.414 9.286-9.286A7.93 7.93 0 1 0 15.848 0" class="primary" clip-rule="evenodd"/></svg>`,"sign-out":Ie`<svg xmlns="http://www.w3.org/2000/svg" class="icon-door-exit" viewBox="0 0 24 24"><path d="M11 4h3a1 1 0 0 1 1 1v3a1 1 0 0 1-2 0V6h-2v12h2v-2a1 1 0 0 1 2 0v3a1 1 0 0 1-1 1h-3v1a1 1 0 0 1-1.27.96l-6.98-2A1 1 0 0 1 2 19V5a1 1 0 0 1 .75-.97l6.98-2A1 1 0 0 1 11 3z" class="primary"/><path d="m18.59 11-1.3-1.3c-.94-.94.47-2.35 1.42-1.4l3 3a1 1 0 0 1 0 1.4l-3 3c-.95.95-2.36-.46-1.42-1.4l1.3-1.3H14a1 1 0 0 1 0-2z" class="secondary"/></svg>`,sort:Ie`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 76 76"><rect width="70" height="11" x="3" y="16" class="primary" rx="5"/><rect width="62" height="11" x="11" y="33" class="primary" rx="5"/><rect width="54" height="11" x="19" y="50" class="primary" rx="5"/></svg>`,trash:Ie`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 57 58"><path stroke-width="2" d="M6.13 18.658h44.74a4 4 0 0 1 3.918 4.804l-6.023 29.356a4 4 0 0 1-4.232 3.184L28.97 54.778a6 6 0 0 0-.94 0l-15.563 1.224a4 4 0 0 1-4.232-3.184L2.212 23.462a4 4 0 0 1 3.918-4.805" class="primary-stroke"/><rect width="4" height="30.964" class="secondary" rx="2" transform="matrix(.99209 -.12553 .2006 .97967 9.295 22.952)"/><rect width="4" height="30.964" class="secondary" rx="2" transform="matrix(.9921 .12548 -.20051 .9797 44.157 22.45)"/><rect width="4" height="28.805" x="26.872" y="22.138" class="secondary" rx="2"/><path fill-rule="evenodd" d="M37.036 0a3.68 3.68 0 0 1 3.678 3.679 3.68 3.68 0 0 0 3.679 3.678h9.664a2.943 2.943 0 0 1 0 5.886H2.943a2.943 2.943 0 0 1 0-5.886h9.664a3.68 3.68 0 0 0 3.679-3.678A3.68 3.68 0 0 1 19.964 0zM22.564 2.207a2.207 2.207 0 1 0 0 4.415h11.872a2.207 2.207 0 0 0 0-4.415z" class="primary" clip-rule="evenodd"/></svg>`,unlink:Ie`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 100 100"><path stroke-width="8" d="m41.035 46-24.49 24.452a5 5 0 0 0 0 7.077l5.957 5.945a5 5 0 0 0 7.065 0L46 67.067M58.195 54l25.103-24.393a5 5 0 0 0-.01-7.183l-6.276-6.06a5 5 0 0 0-6.957.009L53 32.933" class="primary-stroke"/><rect width="8" height="18" x="65" y="74.997" class="shadow" rx="4" transform="rotate(-45 65 74.997)"/><rect width="8" height="18" x="73.498" y="63.489" class="shadow" rx="4" transform="rotate(-75 73.498 63.489)"/><rect width="8" height="18" x="49.681" y="79.357" class="shadow" rx="4" transform="rotate(-15 49.68 79.357)"/><rect width="8" height="18" x="34.445" y="21.543" class="shadow" rx="4" transform="rotate(135 34.445 21.543)"/><rect width="8" height="18" x="24.947" y="33.05" class="shadow" rx="4" transform="rotate(105 24.947 33.05)"/><rect width="8" height="18" x="49.765" y="18.182" class="shadow" rx="4" transform="rotate(165 49.765 18.182)"/></svg>`,"view-hidden":Ie`<svg xmlns="http://www.w3.org/2000/svg" class="icon-view-hidden" viewBox="0 0 24 24"><path d="M15.1 19.34a8 8 0 0 1-8.86-1.68L1.3 12.7a1 1 0 0 1 0-1.42L4.18 8.4l2.8 2.8a5 5 0 0 0 5.73 5.73l2.4 2.4zM8.84 4.6a8 8 0 0 1 8.7 1.74l4.96 4.95a1 1 0 0 1 0 1.42l-2.78 2.78-2.87-2.87a5 5 0 0 0-5.58-5.58L8.85 4.6z" class="primary"/><path d="m3.3 4.7 16 16a1 1 0 0 0 1.4-1.4l-16-16a1 1 0 0 0-1.4 1.4" class="secondary"/></svg>`,"view-visible":Ie`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M17.56 17.66a8 8 0 0 1-11.32 0L1.3 12.7a1 1 0 0 1 0-1.42l4.95-4.95a8 8 0 0 1 11.32 0l4.95 4.95a1 1 0 0 1 0 1.42l-4.95 4.95zM11.9 17a5 5 0 1 0 0-10 5 5 0 0 0 0 10" class="primary"/><circle cx="12" cy="12" r="3" class="secondary"/></svg>`,xmark:Ie`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 33 33"><circle cx="16.149" cy="16.149" r="16.149" class="primary"/><path stroke="#fff" stroke-width="3" d="m9.81 9.96 6.34 6.34m6.339 6.339-6.34-6.339m0 0 6.34-6.34m-6.34 6.34-6.338 6.339"/></svg>`};var qt,Jt=function(e,t,r,i){var o,n=arguments.length,a=n<3?t:null===i?i=Object.getOwnPropertyDescriptor(t,r):i;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(e,t,r,i);else for(var s=e.length-1;s>=0;s--)(o=e[s])&&(a=(n<3?o(a):n>3?o(t,r,a):o(t,r))||a);return n>3&&a&&Object.defineProperty(t,r,a),a};let Qt=qt=class extends tt{constructor(){super(...arguments),this.name="info",this.size="24px",this.hoverable=!1,this.colorway="primary"}render(){var e;const t=null!==(e=qt.colorways[this.colorway])&&void 0!==e?e:qt.colorways.primary;return Ie`
      <div
        class=${zt({hoverable:this.hoverable})}
        style="
          --size: ${this.size};
          --primary: ${t.primary};
          --secondary: ${t.secondary};
          --shadow: ${t.shadow};
        "
      >
        ${Gt[this.name]}
      </div>
    `}};function Zt(e){return"function"==typeof e?e():e}Qt.colorways={primary:{primary:"var(--primary-600)",secondary:"var(--primary-500, #327eff)",shadow:"var(--gray-400, #989898)"},danger:{primary:"var(--danger-600, red)",secondary:"var(--danger-500, pink)",shadow:"var(--gray-500, #888)"}},Qt.styles=re`
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
  `,Jt([at()],Qt.prototype,"name",void 0),Jt([at()],Qt.prototype,"size",void 0),Jt([at({type:Boolean})],Qt.prototype,"hoverable",void 0),Jt([at()],Qt.prototype,"colorway",void 0),Qt=qt=Jt([it("ui-icon")],Qt);class Xt extends Event{static{this.eventName="lit-state-changed"}constructor(e,t,r){super(Xt.eventName,{cancelable:!1}),this.key=e,this.value=t,this.state=r}}const er=(e,t)=>t!==e&&(t==t||e==e);class tr extends EventTarget{static{this.finalized=!1}static initPropertyMap(){this.propertyMap||(this.propertyMap=new Map)}get propertyMap(){return this.constructor.propertyMap}get stateValue(){return Object.fromEntries([...this.propertyMap].map((([e])=>[e,this[e]])))}constructor(){super(),this.hookMap=new Map,this.constructor.finalize(),this.propertyMap&&[...this.propertyMap].forEach((([e,t])=>{if(void 0!==t.initialValue){const r=Zt(t.initialValue);this[e]=r,t.value=r}}))}static finalize(){if(this.finalized)return!1;this.finalized=!0;const e=Object.keys(this.properties||{});for(const t of e)this.createProperty(t,this.properties[t]);return!0}static createProperty(e,t){this.finalize();const r="symbol"==typeof e?Symbol():`__${e}`,i=this.getPropertyDescriptor(String(e),r,t);Object.defineProperty(this.prototype,e,i)}static getPropertyDescriptor(e,t,r){const i=r?.hasChanged||er;return{get(){return this[t]},set(r){const o=this[e];this[t]=r,!0===i(r,o)&&this.dispatchStateEvent(e,r,this)},configurable:!0,enumerable:!0}}reset(){this.hookMap.forEach((e=>e.reset())),[...this.propertyMap].filter((([e,t])=>!(!0===t.skipReset||void 0===t.resetValue))).forEach((([e,t])=>{this[e]=t.resetValue}))}subscribe(e,t,r){t&&!Array.isArray(t)&&(t=[t]);const i=r=>{t&&!t.includes(r.key)||e(r.key,r.value,this)};return this.addEventListener(Xt.eventName,i,r),()=>this.removeEventListener(Xt.eventName,i)}dispatchStateEvent(e,t,r){this.dispatchEvent(new Xt(e,t,r))}}class rr{constructor(e,t,r){this.host=e,this.state=t,this.callback=r||(()=>this.host.requestUpdate()),this.host.addController(this)}hostConnected(){this.state.addEventListener(Xt.eventName,this.callback),this.callback()}hostDisconnected(){this.state.removeEventListener(Xt.eventName,this.callback)}}function ir(e){return(t,r)=>{if(Object.getOwnPropertyDescriptor(t,r))throw new Error("@property must be called before all state decorators");const i=t.constructor;i.initPropertyMap();const o=t.hasOwnProperty(r);return i.propertyMap.set(r,{...e,initialValue:e?.value,resetValue:e?.value}),i.createProperty(r,e),o?Object.getOwnPropertyDescriptor(t,r):void 0}}new URL(location.href);const or={prefix:"_ls"};function nr(e){return e={...or,...e},(t,r)=>{const i=Object.getOwnPropertyDescriptor(t,r);if(!i)throw new Error("@local-storage decorator need to be called after @property");const o=`${e?.prefix||""}_${e?.key||String(r)}`,n=t.constructor,a=n.propertyMap.get(r),s=a?.type;if(a){const t=a.initialValue;a.initialValue=()=>function(e,t){if(null!==e&&(t===Boolean||t===Number||t===Array||t===Object))try{e=JSON.parse(e)}catch(t){console.warn("cannot parse value",e)}return e}(localStorage.getItem(o),s)??Zt(t),n.propertyMap.set(r,{...a,...e})}const l=i?.set,c={...i,set:function(e){void 0!==e&&localStorage.setItem(o,s===Object||s===Array?JSON.stringify(e):e),l&&l.call(this,e)}};Object.defineProperty(n.prototype,r,c)}}var ar=function(e,t,r,i){var o,n=arguments.length,a=n<3?t:null===i?i=Object.getOwnPropertyDescriptor(t,r):i;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(e,t,r,i);else for(var s=e.length-1;s>=0;s--)(o=e[s])&&(a=(n<3?o(a):n>3?o(t,r,a):o(t,r))||a);return n>3&&a&&Object.defineProperty(t,r,a),a};class sr extends tr{get info(){return this.infoRaw?JSON.parse(this.infoRaw):void 0}set info(e){this.infoRaw=e?JSON.stringify(e):""}get username(){var e,t;return null!==(t=null===(e=this.info)||void 0===e?void 0:e.username.split("@")[0])&&void 0!==t?t:"unknown"}get email(){var e,t;return null!==(t=null===(e=this.info)||void 0===e?void 0:e.email)&&void 0!==t?t:"unknown"}get role(){var e,t;return null!==(t=null===(e=this.info)||void 0===e?void 0:e.role)&&void 0!==t?t:"unknown"}set extras(e){this.rawExtras=JSON.stringify(e)}get extras(){if(""===this.rawExtras||void 0===this.rawExtras)return{};try{return JSON.parse(this.rawExtras)}catch(e){return console.error(e),{}}}async loadIfNeeded(e=!1){var t;const r=Date.now(),i=""!==this.infoRaw,o=(null===(t=this.permissions)||void 0===t?void 0:t.length)>0,n=0!=+this.loadedAt&&r-+this.loadedAt>9e5;if(!i||!o||n||e)try{const e=await fetch("/api/management/me");if(!e.ok)throw new Error("Failed to fetch /me");const t=await e.json();if(this.info=t.info,this.permissions=t.permissions.join(","),this.loadedAt=`${r}`,sr.GetExtraAboutMe){const e=await sr.GetExtraAboutMe({id:t.info.id,username:t.info.username,role:t.info.role,permissions:t.info.permissions});this.extras=e}}catch(e){console.error("Failed to load /me:",e),window.location.href=`/login?b=${encodeURIComponent(window.location.pathname+window.location.search)}&utm_source=locksmith&utm_campaign=session_expired`}}hasPermission(e){return(this.permissions||"").split(",").includes(e)}hasRole(e){return Array.isArray(e)?e.some((e=>this.role===e)):this.role===e}get isLaunchpad(){const e=document.cookie.split(";");for(let t of e)if(t=t.trim(),t.startsWith("LaunchpadUser="))return t.substring(14)}clear(){localStorage.removeItem("_identity_i"),localStorage.removeItem("_identity_p"),localStorage.removeItem("_identity_la"),localStorage.removeItem("_identity_x")}signOut(){void 0!==sr.SignOutCallback&&sr.SignOutCallback({id:this.info.id,username:this.info.username,role:this.info.role,permissions:this.permissions.split(",")}),this.clear(),window.location.href="/sign-out"}}sr.SignOutCallback=void 0,sr.GetExtraAboutMe=void 0,ar([nr({key:"i",prefix:"_identity"}),ir()],sr.prototype,"infoRaw",void 0),ar([nr({key:"x",prefix:"_identity"}),ir()],sr.prototype,"rawExtras",void 0),ar([nr({key:"p",prefix:"_identity"}),ir()],sr.prototype,"permissions",void 0),ar([nr({key:"la",prefix:"_identity"}),ir()],sr.prototype,"loadedAt",void 0);const lr=new sr;var cr=function(e,t,r,i){var o,n=arguments.length,a=n<3?t:null===i?i=Object.getOwnPropertyDescriptor(t,r):i;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(e,t,r,i);else for(var s=e.length-1;s>=0;s--)(o=e[s])&&(a=(n<3?o(a):n>3?o(t,r,a):o(t,r))||a);return n>3&&a&&Object.defineProperty(t,r,a),a};let dr=class extends tt{constructor(){super(...arguments),this.location={vertical:"top",horizontal:"right"},this.open=!1,this.launchpadForceClosed=!1,this.aboutMeState=new rr(this,lr),this.manageAccountButton=Ct(),this.handleEscape=e=>{"Escape"===e.key&&(this.open=!1)},this.listenForOOBClicks=e=>{if(!this.shadowRoot)return;e.composedPath().includes(this.shadowRoot.host)||(this.open=!1,window.removeEventListener("click",this.listenForOOBClicks))}}updated(){this.setAttribute("location-vertical",this.location.vertical),this.setAttribute("location-horizontal",this.location.horizontal)}connectedCallback(){super.connectedCallback(),lr.loadIfNeeded()}openClicked(){this.open=!this.open,setTimeout((()=>{this.open?(window.addEventListener("click",this.listenForOOBClicks),window.addEventListener("keydown",this.handleEscape),this.updateComplete.then((()=>{setTimeout((()=>{var e;null===(e=this.manageAccountButton.value)||void 0===e||e.focus()}),100)}))):(window.removeEventListener("click",this.listenForOOBClicks),window.removeEventListener("keydown",this.handleEscape))}))}render(){return Ie` <div id="container">
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
              <p id="title">${lr.username}</p>
              <p id="desc">${lr.email}</p>
            </div>
          </div>

          <div id="actions">
            <button
              ${Rt(this.manageAccountButton)}
              @click=${()=>{window.location.href="/profile"}}
            >
              <ui-icon name="cog" size="1rem"></ui-icon>
              Manage Account
              <span></span>
            </button>
            ${lr.hasPermission("view.ls-admin")?Ie`
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
              @click=${()=>{lr.signOut()}}
            >
              <ui-icon name="sign-out" size="1rem" colorway="danger"></ui-icon>
              Sign out
            </button>
          </div>
        </div>
      </div>

      ${void 0===lr.isLaunchpad||this.launchpadForceClosed?void 0:Ie`<div id="launchpad">
            <div>
              <p id="launchpad-status">
                ${lr.isLaunchpad.toUpperCase()} Role
              </p>
              <p id="launchpad-user">
                This is what the ${lr.isLaunchpad.toUpperCase()} role would
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
          </div>`}`}};dr.styles=[re`
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
    `],cr([at()],dr.prototype,"location",void 0),cr([st()],dr.prototype,"open",void 0),cr([st()],dr.prototype,"launchpadForceClosed",void 0),dr=cr([it("locksmith-user-icon")],dr);var pr=function(e,t,r,i){var o,n=arguments.length,a=n<3?t:null===i?i=Object.getOwnPropertyDescriptor(t,r):i;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(e,t,r,i);else for(var s=e.length-1;s>=0;s--)(o=e[s])&&(a=(n<3?o(a):n>3?o(t,r,a):o(t,r))||a);return n>3&&a&&Object.defineProperty(t,r,a),a};let hr=0,ur=class extends tt{constructor(){super(),this.toasts=[],window.addEventListener("close-toast",(e=>{const t=e.detail.id;this.removeToastById(t)})),window.addEventListener("do-toast",(e=>{const{id:t,text:r,danger:i,persist:o,actionText:n,onClick:a,duration:s}=e.detail,l=void 0!==s?s:i?8e3:2500,c=null!=t?t:hr++,d={id:c,text:r,danger:i,duration:l,persist:o,actionText:n,onClick:a};this.toasts=[...this.toasts,d],o||setTimeout((()=>this.removeToastById(c)),l)}))}removeToastById(e){const t=this.toasts.find((t=>t.id===e));t&&(t.removing=!0,this.toasts=[...this.toasts],setTimeout((()=>{this.toasts=this.toasts.filter((t=>t.id!==e))}),200))}render(){return Ie`<div id="root">
      ${St(this.toasts,(e=>e.id),(e=>{var t;return Ie`
          <div
            class=${zt({toast:!0,danger:e.danger,"slide-out":null!==(t=e.removing)&&void 0!==t&&t})}
            style="--toast-duration: ${e.duration}ms"
          >
            ${e.text}
            ${e.actionText&&e.onClick?Ie`<button
                  class="action-btn"
                  @click=${()=>{var t;null===(t=e.onClick)||void 0===t||t.call(e),this.removeToastById(e.id)}}
                >
                  ${e.actionText}
                </button>`:null}
            ${e.persist&&void 0===e.onClick?Ie`<button
                  class="close-btn"
                  @click=${()=>this.removeToastById(e.id)}
                >
                  &times;
                </button>`:null}

            <div
              class="progress-bar"
              style=${e.persist?"display: none":`animation: shrink ${e.duration}ms linear forwards`}
            ></div>
          </div>
        `}))}
    </div>`}};function mr(e){var t,r;window.dispatchEvent(new CustomEvent("do-toast",{detail:{id:e.id,text:e.text,danger:null!==(t=e.danger)&&void 0!==t&&t,persist:null!==(r=e.persist)&&void 0!==r&&r,actionText:e.actionText,onClick:e.onClick,duration:e.duration}}))}ur.styles=[re`
      #root {
        z-index: 100000000;
        position: fixed;
        bottom: 2rem;
        left: 0;
        display: flex;
        flex-direction: column-reverse;
        align-items: center;
        width: 100%;
        gap: 0.5rem;
        box-sizing: border-box;
      }

      .toast {
        box-sizing: border-box;
        --progress-color: var(--primary-300);
        --progress-background: var(--primary-100);
        --toast-height: 0.25rem;
        width: calc(100% - 4rem);
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 0.25rem;
        padding: 0.5rem 1rem;
        border-radius: 0.25rem;
        background-color: #fff;
        border: 1px solid var(--primary-200);
        color: var(--primary-800);
        box-shadow: rgba(0, 0, 0, 0.05) 0px 3px 8px;
        position: relative;
        overflow: hidden;
        padding-bottom: calc(var(--toast-height) + 0.5rem);
        animation: slideIn 200ms ease-out;
      }

      .progress-bar {
        position: absolute;
        bottom: 0;
        left: 0;
        height: var(--toast-height);
        width: 100%;
        background-color: var(--progress-color);
      }

      .toast::before {
        content: " ";
        position: absolute;
        bottom: 0;
        left: 0;
        height: var(--toast-height);
        width: 100%;
        background-color: var(--progress-background);
      }

      .danger {
        --progress-background: var(--danger-200);
        --progress-color: var(--danger-400);
        background-color: var(--danger-50);
        border: 1px solid var(--danger-200);
        color: var(--danger-800);
      }

      .toast.slide-out {
        animation: slideOut 200ms ease-in forwards;
      }

      .close-btn {
        background: transparent;
        border: none;
        font-weight: bold;
        cursor: pointer;
        color: inherit;
        font-size: 1rem;
      }

      .action-btn {
        background: none;
        border: none;
        color: var(--primary-600);
        font-weight: 600;
        cursor: pointer;
        text-decoration: underline;
      }

      .toast.danger .action-btn {
        color: var(--danger-600);
      }

      @media (min-width: 680px) {
        .toast {
          width: fit-content;
          min-width: 12rem;
          justify-content: center;
        }
      }

      @keyframes slideOut {
        from {
          transform: translateY(0);
          opacity: 1;
        }
        to {
          transform: translateY(2rem);
          opacity: 0;
        }
      }

      @keyframes shrink {
        from {
          width: 100%;
        }
        to {
          width: 0%;
        }
      }

      @keyframes slideIn {
        from {
          transform: translateY(2rem);
          opacity: 0;
        }
        to {
          transform: translateY(0);
          opacity: 1;
        }
      }
    `],pr([st()],ur.prototype,"toasts",void 0),ur=pr([it("toast-component")],ur);var fr=function(e,t,r,i){var o,n=arguments.length,a=n<3?t:null===i?i=Object.getOwnPropertyDescriptor(t,r):i;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(e,t,r,i);else for(var s=e.length-1;s>=0;s--)(o=e[s])&&(a=(n<3?o(a):n>3?o(t,r,a):o(t,r))||a);return n>3&&a&&Object.defineProperty(t,r,a),a};let vr=class extends tt{constructor(){super(),this.locale="en",this.mobileNavOpen=!1,this.includeFooter=!0,this.navbars={};this.navbars={user:{left:[{name:"Overview",path:"/locksmith"},{name:"Users",path:"/locksmith/users",dropdown:[]}],right:[]}},window.ononline=this.onlineNotice,window.onoffline=this.offlineNotice}onlineNotice(){var e;console.log("notyfing toast"),e="network",window.dispatchEvent(new CustomEvent("close-toast",{detail:{id:e}})),mr({id:"net",text:"You are back online."})}offlineNotice(){mr({id:"network",text:"You are offline. Please reconnect to WiFi.",persist:!0,danger:!0})}removeTrailingSlash(e){return e.endsWith("/")?e.slice(0,-1):e}renderDropdownNavbar(e){return Ie`<section class="dropdown">
      <a
        class="${this.removeTrailingSlash(window.location.pathname)===e.path?"active":""}"
        href="${e.path}"
        >${e.name}</a
      >
      ${void 0!==e.dropdown?Ie`
            <section>
              ${e.dropdown.map((e=>Ie`<button @click=${()=>J.go(e.path)}>
                    ${e.name}
                  </button>`))}
            </section>
          `:""}
    </section>`}render(){return Ie`<div id="root">
        <prompt-component></prompt-component>
        <alert-component></alert-component>
        <toast-component></toast-component>
      <div id="navcontainer">
        <nav>
          <div id="full">
            <div>
              <h1 id="logo">
                  <svg viewBox="0 0 46 50" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M23 20V27.5C21.901 27.5019 20.8332 27.8659 19.9618 28.5357C19.0904 29.2055 18.464 30.1437 18.1794 31.2052C17.8948 32.2668 17.9679 33.3925 18.3874 34.4084C18.8068 35.4242 19.5493 36.2735 20.5 36.825V40C20.5 40.663 20.7634 41.2989 21.2322 41.7678C21.7011 42.2366 22.337 42.5 23 42.5V50H5.5C4.17392 50 2.90215 49.4732 1.96447 48.5355C1.02678 47.5979 0.5 46.3261 0.5 45V25C0.5 22.25 2.75 20 5.5 20H23Z" fill="#173ab6"/>
                  <path d="M23 42.5C23.663 42.5 24.2989 42.2366 24.7678 41.7678C25.2366 41.2989 25.5 40.663 25.5 40V36.825C26.4507 36.2735 27.1932 35.4242 27.6126 34.4084C28.0321 33.3925 28.1052 32.2668 27.8206 31.2052C27.536 30.1437 26.9096 29.2055 26.0382 28.5357C25.1668 27.8659 24.099 27.5019 23 27.5V20H30.5V12.5C30.5 10.5109 29.7098 8.60322 28.3033 7.1967C26.8968 5.79018 24.9891 5 23 5C21.0109 5 19.1032 5.79018 17.6967 7.1967C16.2902 8.60322 15.5 10.5109 15.5 12.5V15H10.5V12.5C10.5 9.18479 11.817 6.00537 14.1612 3.66117C16.5054 1.31696 19.6848 0 23 0C26.3152 0 29.4946 1.31696 31.8388 3.66117C34.183 6.00537 35.5 9.18479 35.5 12.5V20H40.5C41.8261 20 43.0979 20.5268 44.0355 21.4645C44.9732 22.4021 45.5 23.6739 45.5 25V45C45.5 46.3261 44.9732 47.5979 44.0355 48.5355C43.0979 49.4732 41.8261 50 40.5 50H23V42.5Z" fill="#173ab6"/>
                  </svg>
              </h1>
              ${this.navbars.user.left.map((e=>void 0!==e.dropdown?this.renderDropdownNavbar(e):Ie`<a
                  class="${this.removeTrailingSlash(window.location.pathname)===this.removeTrailingSlash(e.path)?"active":""}"
                  href="${e.path}"
                  >${e.name}</a
                >`))}
            </div>

            <div>
              <!-- <global-search-component></global-search-component> -->

              ${this.navbars.user.right.map((e=>void 0!==e.dropdown?this.renderDropdownNavbar(e):Ie`<a
                  class="${this.removeTrailingSlash(window.location.pathname)===this.removeTrailingSlash(e.path)?"active":""}"
                  href="${e.path}"
                  >${e.name}</a
                >`))}
              <locksmith-user-icon .location=${{vertical:"bottom",horizontal:"right"}}></locksmith-user-icon>

              <button
                id="mobile-ham"
                @click=${()=>{this.mobileNavOpen=!this.mobileNavOpen}}
              >
                &#x2630;
              </button>
            </div>
          </div>
          <div id="mobile" class="${this.mobileNavOpen?"show":""}">
            ${this.navbars.user.left.map((e=>Ie`<a
                  class="${window.location.pathname===e.path?"active":""}"
                  href="${e.path}"
                  >${e.name}</a
                >`))}
            ${this.navbars.user.right.map((e=>Ie`<a
                  class="${window.location.pathname===e.path?"active":""}"
                  href="${e.path}"
                  >${e.name}</a
                >`))}
            <a
              href="/dashboard/settings"
              aria-label="Settings"
              class="${"/dashboard/settings"===window.location.pathname?"active":""}"
              ><span class="icon cog md settings"></span> Settings</a
            >
          </div>
        </nav>
      </div>
      <div id="wrapper">
        <slot name="aside"><br></br></slot>
        <main>
          <slot></slot>
        </main>
        <slot name="aside-right"><br></br></slot>
      </div>
    </div>`}};vr.styles=[lt,re`
      :host {
        display: block;
        overflow-x: hidden;
      }
      #root {
        min-height: 100dvh;
        height: 100%;
        display: flex;
        align-items: center;
        flex-direction: column;
        position: relative;
        background-color: #f3f3fa;
        width: 100vw;
      }

      #navcontainer {
        background-color: rgb(255, 255, 255);
        padding: 0.5rem 1rem;
        position: fixed;
        top: 0px;
        left: 0px;
        width: 100vw;
        gap: 1.25rem;
        display: flex;
        justify-content: center;
        border-bottom: 1px solid var(--gray-200);
        z-index: 1000;
      }

      h1#logo {
        display: flex;
        align-items: center;
        justify-content: center;
        height: 1rem;
        aspect-ratio: 1/1;
      }

      nav {
        width: 100%;
        max-width: 72rem;
      }

      nav #full {
        width: 100%;
        display: flex;
        gap: 1rem;
        align-items: center;
        flex-direction: row;
        justify-content: space-between;
        padding: 0.35rem 0rem;
      }

      nav #full div:not(#overview) {
        display: flex;
        align-items: center;
        gap: 1.25rem;
      }

      nav p,
      nav a,
      nav button {
        padding: 0.45rem 1rem;
        border: 0;
      }
      nav a,
      nav button {
        color: var(--primary-950);
        background-color: white;
        border-radius: 0.75rem;
        transition: 200ms;
        font-weight: 500;
        cursor: pointer;
        text-decoration: none;
        display: flex;
        align-items: center;
        gap: 0.75rem;
      }

      nav button {
        font-size: 1rem;
        width: 100%;
        text-align: left;
        font-weight: 500;
        white-space: no-wrap;
      }

      nav img {
        width: 3rem;
        aspect-ratio: 1 / 1;
      }

      nav #full a:not(.active),
      nav #full p,
      nav #full span {
        display: none;
      }

      nav a.active,
      nav button.active {
        background-color: var(--primary-100);
        color: var(--primary-800);
      }

      nav span.settings {
        cursor: pointer;
      }

      nav a:not(.active):hover,
      nav button:not(.active):hover {
        background-color: var(--primary-50);
      }

      #mobile {
        flex-direction: column;
        gap: 0.5rem;
        padding-top: 1rem;
        display: none;
      }

      #mobile.show {
        display: flex;
      }

      #mobile a,
      #mobile p {
        padding: 0.75rem 1rem 0.75rem 1rem;
      }

      #mobile #footer {
        display: flex;
        justify-content: space-between;
      }

      #wrapper {
        margin-top: 69px;
        display: grid;
        gap: 2rem;
        padding: 1rem;
        width: 100%;
      }

      h1 {
        font-size: 1rem;
        color: #6c43e8;
      }

      .dropdown {
        position: relative;
      }

      .dropdown > section {
        display: none;
        position: absolute;
        box-shadow: rgba(0, 0, 0, 0.24) 0px 3px 8px;
        border-radius: 0.75rem;
        padding: 0.25rem;
        z-index: 1001;
        width: 100%;
        background-color: white;
      }

      #mobile-ham {
        background-color: #fff;
        border: 0;
        font-size: 1.5rem;
        color: #3d2a7c;
      }

      nav #full .medium-hide {
        display: none;
      }

      #footer {
        color: var(--gray-400);
        font-size: 0.75rem;
        margin-top: 2rem;
        padding-top: 2rem;
        border-top: 1px solid var(--gray-200);
      }

      .desktop-hide {
        display: inline-block !important;
      }

      #overview {
        display: none;
      }

      @media (min-width: 680px) {
        .dropdown:hover > section {
          display: block;
        }
        .desktop-hide {
          display: none !important;
        }
        nav img {
          width: 3rem;
        }

        nav #full a:not(.active),
        nav #full p,
        nav #full span {
          display: block;
        }

        #mobile-ham,
        #mobile,
        #mobile.show {
          display: none;
        }
      }

      #saving {
        background-color: var(--primary-200);
        border-radius: 2rem;
        color: var(--primary-800);
      }

      @media (min-width: 825px) {
        nav #full .medium-hide {
          display: block;
        }

        :host([wide-content]) main {
          max-width: 65vw;
          max-width: 65dvw;
        }
      }

      @media (min-width: 930px) {
        #wrapper {
          margin-top: 80px;
          padding: 1rem 0;
          grid-template-columns: 1fr minmax(0, 72rem) 1fr;
        }
      }
    `],fr([at({type:String})],vr.prototype,"locale",void 0),fr([at({type:Boolean})],vr.prototype,"mobileNavOpen",void 0),fr([at({type:Boolean})],vr.prototype,"includeFooter",void 0),fr([at({type:Object})],vr.prototype,"navbars",void 0),vr=fr([it("locksmith-layout")],vr);var gr=function(e,t,r,i){var o,n=arguments.length,a=n<3?t:null===i?i=Object.getOwnPropertyDescriptor(t,r):i;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(e,t,r,i);else for(var s=e.length-1;s>=0;s--)(o=e[s])&&(a=(n<3?o(a):n>3?o(t,r,a):o(t,r))||a);return n>3&&a&&Object.defineProperty(t,r,a),a};let yr=class extends tt{constructor(){super(),this.urlBase="",this.urlGroup="",this.urlKey="",this.pages={},this.activePageKey="",this.activePageComponent=void 0,this.activePage=void 0,this.listenForPageChanges=e=>{const{location:t}=e.detail,{pathname:r}=t;this.destruct(),0!==r.substring(this.urlBase.length,r.length).length&&(this.urlKey=t.params.urlKey,this.urlGroup=t.params.urlGroup)},window.addEventListener("vaadin-router-location-changed",this.listenForPageChanges),this.addEventListener("destruct",this.destruct)}destruct(){window.removeEventListener("vaadin-router-location-changed",this.listenForPageChanges),this.removeEventListener("destruct",this.destruct)}firstUpdated(){setTimeout((()=>{const e=this.getRouteParams();if(this.urlGroup=e.urlGroup||"",this.urlKey=e.urlKey||"",""===this.activePageKey){if(""===this.urlGroup||""===this.urlKey)return void this.loadPage(Object.keys(this.pages)[0],0);const e=Object.keys(this.pages).find((e=>e.toLowerCase()===this.urlGroup.toString().toLowerCase()));if(e){const t=this.pages[e].find((e=>e.PageKey===this.urlKey));t?this.loadPage(e,t.SortIndex):this.loadPage(Object.keys(this.pages)[0],0)}else this.loadPage(Object.keys(this.pages)[0],0)}}),0)}getRouteParams(){if(Pr){return Pr.location.params}return{}}async loadPage(e,t,r){this.activePageComponent=void 0;const i=this.pages[e].findIndex((e=>e.SortIndex===t));if(-1===i)throw new Error("page index does not exist!");const o=this.pages[e][i];if(this.urlGroup!==e.toLowerCase()&&this.urlKey!==o.PageKey||o.PageKey!==this.activePageKey){const t=`${this.urlBase}/${e.toLowerCase()}/${o.PageKey}`;r||(""===this.activePageKey?window.history.replaceState({group:e,page:o.PageKey},"",t):window.history.pushState({group:e,page:o.PageKey},"",t))}this.activePage=o,this.activePageKey=o.PageKey;const n=new o.PageComponent;if(void 0!==o.LoadProps){const e=o.LoadProps();n.setProps(e)}await n.OnPageLoad(),this.activePageComponent=n}render(){var e,t;return Ie` <locksmith-layout>
      <aside class="aside" slot="aside">
        ${Object.keys(this.pages).map((e=>Ie`<div>
              <h3>${e}</h3>
              <div id="buttons">
                ${this.pages[e].sort(((e,t)=>e.SortIndex-t.SortIndex)).map((t=>Ie`<button
                        @click=${()=>this.loadPage(e,t.SortIndex)}
                        class="${this.activePageKey===t.PageKey?"active":""}"
                      >
                        ${t.PageName}
                      </button>`))}
              </div>
            </div>`))}
        ${void 0!==(null===(e=this.activePageComponent)||void 0===e?void 0:e.ExtraLeftAside)?Ie`<hr />
              ${this.activePageComponent.ExtraLeftAside()}`:Ie``}
      </aside>

      <div id="main">
        ${void 0!==this.activePageComponent?Ie`${this.activePageComponent}`:"Loading.."}
      </div>

      ${void 0!==(null===(t=this.activePageComponent)||void 0===t?void 0:t.GetRightAside)?Ie`<div class="aside right" slot="aside-right">
            ${this.activePageComponent.GetRightAside()}
          </div>`:Ie``}
    </locksmith-layout>`}};yr.styles=[lt,re`
      #root {
      }

      .aside div#buttons {
        display: flex;
        flex-wrap: wrap;
        gap: 1rem;
        width: 100%;
        row-gap: 0.5rem;
      }

      .aside > div:first-of-type:last-of-type > h3 {
        display: none;
      }

      .aside > div {
        display: flex;
        flex-direction: column;
        gap: 0.15rem;
      }
      .aside > div:not(:last-of-type) {
        margin-bottom: 1rem;
      }
      .aside div h3 {
        font-size: 0.85rem;
        color: var(--primary-800);
        margin-bottom: 0.5rem;
      }

      .aside button {
        background-color: var(--primary-200);
        border: 0;
        padding: 0.5rem 1rem;
        border-radius: 2rem;
        font-size: 1rem;
        color: var(--primary-950);
      }

      .aside button.active {
        background-color: var(--primary-500);
        color: #fff;
      }

      @media (min-width: 825px) {
        .aside {
          justify-self: flex-end;
          display: flex;
          flex-direction: column;
          gap: 0.5rem;
          width: 100%;
          width: 14rem;
        }

        .aside:not(.right) {
          padding-left: 2rem;
        }

        .aside > div {
          display: flex;
          flex-direction: column;
        }
        .aside div#buttons {
          width: 100%;
          display: flex;
          flex-direction: column;
          gap: 0.5rem;
        }

        .aside div h3 {
          font-size: 1rem;
          margin-bottom: 1rem;
          color: black;
        }

        .aside > div:first-of-type:last-of-type > h3 {
          display: inherit;
        }

        .aside div:not(:first-of-type) {
          margin-top: 1.25rem;
        }

        .aside button {
          width: 100%;
          padding: 1rem 0rem;
          font-size: 1rem;
          background-color: unset;
          color: var(--primary-950);
          text-align: left;
          border: 0;
          border-radius: 0.25rem;
          cursor: pointer;
          transition: 200ms;
        }

        @media (hover) {
          .aside button:hover {
            background-color: var(--primary-100);
            padding: 1rem;
          }
        }

        .aside button.active {
          background-color: var(--primary-200);
          color: var(--primary-950);
          transition: 200ms;
          padding: 1rem;
        }

        .aside button.active:focus-visible {
          outline: 2px solid var(--primary-600);
        }

        #main {
          max-width: 72rem;
        }
      }
    `],gr([at({type:String})],yr.prototype,"urlBase",void 0),gr([at({type:String})],yr.prototype,"urlGroup",void 0),gr([at({type:String})],yr.prototype,"urlKey",void 0),gr([at({type:Array})],yr.prototype,"pages",void 0),gr([st()],yr.prototype,"activePageKey",void 0),gr([st()],yr.prototype,"activePageComponent",void 0),gr([st()],yr.prototype,"activePage",void 0),yr=gr([it("locksmith-subnav-layout")],yr);var br=function(e,t,r,i){var o,n=arguments.length,a=n<3?t:null===i?i=Object.getOwnPropertyDescriptor(t,r):i;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(e,t,r,i);else for(var s=e.length-1;s>=0;s--)(o=e[s])&&(a=(n<3?o(a):n>3?o(t,r,a):o(t,r))||a);return n>3&&a&&Object.defineProperty(t,r,a),a};let wr=class extends tt{constructor(){super(...arguments),this.users=[],this.loading=!0}OnPageLoad(){this.fetchUsers()}async fetchUsers(){const e=await fetch("/api/users/list"),t=await e.json();this.users=t,this.loading=!1}render(){return Ie`
      <header>
        <div
          style="display: flex; flex-direction: row; align-items: center; gap: 1rem; margin-bottom: 1rem;"
        >
          <ui-icon
            name="person-group"
            colorway="lightSecondary"
            size="1.45rem"
          ></ui-icon>

          <h1 style="margin: 0;">All Users</h1>
        </div>
      </header>

      <div class="widget" style="${this.loading?"display: none;":""}">
        <div class="list">
          ${this.users.map((e=>{var t;return Ie`
              <button href="#" class="item">
                <p id="title">${null!==(t=e.username)&&void 0!==t?t:e.email}</p>
                <p id="description">
                  ${e.role} &bull; ${e.sessions} active sessions
                </p>
              </button>
            `}))}
        </div>
      </div>
    `}};async function xr(e){const t=crypto.randomUUID();return e.id=t,new Promise((t=>{window.addEventListener(`fl-response-${e.id}`,(e=>{t(e.detail)}),{once:!0}),window.dispatchEvent(new CustomEvent("fl-prompt",{detail:e}))}))}wr.styles=[lt,re``],br([st()],wr.prototype,"users",void 0),br([st()],wr.prototype,"loading",void 0),wr=br([it("locksmith-users-subpage")],wr);var $r=function(e,t,r,i){var o,n=arguments.length,a=n<3?t:null===i?i=Object.getOwnPropertyDescriptor(t,r):i;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(e,t,r,i);else for(var s=e.length-1;s>=0;s--)(o=e[s])&&(a=(n<3?o(a):n>3?o(t,r,a):o(t,r))||a);return n>3&&a&&Object.defineProperty(t,r,a),a};let kr=class extends tt{constructor(){super(...arguments),this.invites=[],this.loading=!0,this.sendInviteBtnRef=Ct()}OnPageLoad(){this.fetchInvites()}async fetchInvites(){const e=await fetch("/api/users/invitations"),t=await e.json();this.invites=t.reverse(),this.loading=!1}async sendInvite(){const e=await xr({id:"send-invites",title:"Send a new Invite",description:"What is the new users email address? You will be able to select their role on the next screen.",type:"text"});if(e.canceled)return;const t=await xr({id:"invite-role",title:"What is this users role?",description:"This will give them specific permissions, be careful!",type:"radio",radioOptions:[{key:"user",title:"User"},{key:"admin",title:"Admin"}]});if(!t.canceled){this.sendInviteBtnRef.value.loading=!0;try{const i=await async function(e,t){return fetch(e,t)}("/api/users/invite",{method:"POST",body:JSON.stringify({email:e.value,role:t.value})});if(200!==i.status){if(409===i.status)return r={title:"This user has already been invited.",description:"Please use a different email or re-issue the old invite."},void window.dispatchEvent(new CustomEvent("fl-alert",{detail:r}));throw new Error("Got a bad status code: "+i.status)}mr({text:"Invitation email has been sent."}),this.invites=[{email:e.value,role:t.value,inviter:"",sentAt:+new Date/1e3,userid:""},...this.invites]}catch(e){mr({text:`Failed to send invite: ${e.message}`,danger:!0})}finally{this.sendInviteBtnRef.value.loading=!1}var r}}render(){return Ie`
      <header class="x">
        <div
          style="display: flex; flex-direction: row; align-items: center; gap: 1rem; margin-bottom: 1rem;"
        >
          <ui-icon
            name="email"
            colorway="lightSecondary"
            size="1.45rem"
          ></ui-icon>

          <h1 style="margin: 0;">Pending Invitations</h1>
        </div>
        <div>
          <button-component
            .expectLoad=${!0}
            .loadingText=${"Sending Invite.."}
            ${Rt(this.sendInviteBtnRef)}
            @fl-click=${()=>{this.sendInvite()}}
          >
            Send Invite
          </button-component>
        </div>
      </header>

      <div class="widget" style="${this.loading?"display: none;":""}">
        <div class="list">
          ${this.invites.map((e=>Ie`
              <button href="#" class="item">
                <p id="title">${e.email}</p>
                <p id="description">Invited ${function(e){const t=new Date(1e3*e),r=new Date,i=Math.floor((r.getTime()-t.getTime())/1e3);if(i<60)return`${i}s ago`;const o=Math.floor(i/60);if(o<60)return`${o}m ago`;const n=Math.floor(o/60);if(n<24)return`${n}h ago`;const a=Math.floor(n/24);return a<7?`${a}d ago`:t.toLocaleDateString()}(e.sentAt)}</p>
              </button>
            `))}
        </div>
      </div>
    `}};kr.styles=[lt,re``],$r([st()],kr.prototype,"invites",void 0),$r([st()],kr.prototype,"loading",void 0),kr=$r([it("locksmith-invitations-subpage")],kr);var _r=function(e,t,r,i){var o,n=arguments.length,a=n<3?t:null===i?i=Object.getOwnPropertyDescriptor(t,r):i;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(e,t,r,i);else for(var s=e.length-1;s>=0;s--)(o=e[s])&&(a=(n<3?o(a):n>3?o(t,r,a):o(t,r))||a);return n>3&&a&&Object.defineProperty(t,r,a),a};let Ar=class extends tt{constructor(){super(...arguments),this.pages={Users:[{PageKey:"all",PageComponent:wr,PageName:"Registered Users",SortIndex:0},{PageKey:"invites",PageComponent:kr,PageName:"Invitations",SortIndex:1}]}}onBeforeEnter(e){}render(){return Ie`
      <locksmith-subnav-layout
        .urlBase=${"/locksmith/users"}
        .pages=${this.pages}
      >
      </locksmith-subnav-layout>
    `}};Ar.styles=[lt,re``],Ar=_r([it("locksmith-users-page")],Ar),Qt.colorways.primary={primary:"var(--primary-800)",secondary:"var(--primary-300)",shadow:"var(--gray-600)"};const Cr=[{path:"/locksmith/users/:urlGroup?/:urlKey?",component:"locksmith-users-page"},{path:"(.*)",action:e=>{console.log(e),console.warn("Page not found"),J.go("/locksmith/users")},component:"not-found-page"}],Er=document.getElementById("outlet"),Pr=new J(Er);Pr.setRoutes(Cr);export{Pr as router};
//# sourceMappingURL=locksmith-admin.bundle.js.map
