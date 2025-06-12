function t(t,e){void 0===e&&(e={});for(var r=function(t){for(var e=[],r=0;r<t.length;){var i=t[r];if("*"!==i&&"+"!==i&&"?"!==i)if("\\"!==i)if("{"!==i)if("}"!==i)if(":"!==i)if("("!==i)e.push({type:"CHAR",index:r,value:t[r++]});else{var o=1,n="";if("?"===t[s=r+1])throw new TypeError('Pattern cannot start with "?" at '.concat(s));for(;s<t.length;)if("\\"!==t[s]){if(")"===t[s]){if(0==--o){s++;break}}else if("("===t[s]&&(o++,"?"!==t[s+1]))throw new TypeError("Capturing groups are not allowed at ".concat(s));n+=t[s++]}else n+=t[s++]+t[s++];if(o)throw new TypeError("Unbalanced pattern at ".concat(r));if(!n)throw new TypeError("Missing pattern at ".concat(r));e.push({type:"PATTERN",index:r,value:n}),r=s}else{for(var a="",s=r+1;s<t.length;){var l=t.charCodeAt(s);if(!(l>=48&&l<=57||l>=65&&l<=90||l>=97&&l<=122||95===l))break;a+=t[s++]}if(!a)throw new TypeError("Missing parameter name at ".concat(r));e.push({type:"NAME",index:r,value:a}),r=s}else e.push({type:"CLOSE",index:r,value:t[r++]});else e.push({type:"OPEN",index:r,value:t[r++]});else e.push({type:"ESCAPED_CHAR",index:r++,value:t[r++]});else e.push({type:"MODIFIER",index:r,value:t[r++]})}return e.push({type:"END",index:r,value:""}),e}(t),o=e.prefixes,n=void 0===o?"./":o,a=e.delimiter,s=void 0===a?"/#?":a,l=[],c=0,d=0,p="",h=function(t){if(d<r.length&&r[d].type===t)return r[d++].value},u=function(t){var e=h(t);if(void 0!==e)return e;var i=r[d],o=i.type,n=i.index;throw new TypeError("Unexpected ".concat(o," at ").concat(n,", expected ").concat(t))},m=function(){for(var t,e="";t=h("CHAR")||h("ESCAPED_CHAR");)e+=t;return e},f=function(t){var e=l[l.length-1],r=t||(e&&"string"==typeof e?e:"");if(e&&!r)throw new TypeError('Must have text between two parameters, missing text after "'.concat(e.name,'"'));return!r||function(t){for(var e=0,r=s;e<r.length;e++){var i=r[e];if(t.indexOf(i)>-1)return!0}return!1}(r)?"[^".concat(i(s),"]+?"):"(?:(?!".concat(i(r),")[^").concat(i(s),"])+?")};d<r.length;){var v=h("CHAR"),g=h("NAME"),y=h("PATTERN");if(g||y){var b=v||"";-1===n.indexOf(b)&&(p+=b,b=""),p&&(l.push(p),p=""),l.push({name:g||c++,prefix:b,suffix:"",pattern:y||f(b),modifier:h("MODIFIER")||""})}else{var w=v||h("ESCAPED_CHAR");if(w)p+=w;else if(p&&(l.push(p),p=""),h("OPEN")){b=m();var x=h("NAME")||"",$=h("PATTERN")||"",k=m();u("CLOSE"),l.push({name:x||($?c++:""),pattern:x&&!$?f(b):$,prefix:b,suffix:k,modifier:h("MODIFIER")||""})}else u("END")}}return l}function e(e,i){return r(t(e,i),i)}function r(t,e){void 0===e&&(e={});var r=o(e),i=e.encode,n=void 0===i?function(t){return t}:i,a=e.validate,s=void 0===a||a,l=t.map((function(t){if("object"==typeof t)return new RegExp("^(?:".concat(t.pattern,")$"),r)}));return function(e){for(var r="",i=0;i<t.length;i++){var o=t[i];if("string"!=typeof o){var a=e?e[o.name]:void 0,c="?"===o.modifier||"*"===o.modifier,d="*"===o.modifier||"+"===o.modifier;if(Array.isArray(a)){if(!d)throw new TypeError('Expected "'.concat(o.name,'" to not repeat, but got an array'));if(0===a.length){if(c)continue;throw new TypeError('Expected "'.concat(o.name,'" to not be empty'))}for(var p=0;p<a.length;p++){var h=n(a[p],o);if(s&&!l[i].test(h))throw new TypeError('Expected all "'.concat(o.name,'" to match "').concat(o.pattern,'", but got "').concat(h,'"'));r+=o.prefix+h+o.suffix}}else if("string"!=typeof a&&"number"!=typeof a){if(!c){var u=d?"an array":"a string";throw new TypeError('Expected "'.concat(o.name,'" to be ').concat(u))}}else{h=n(String(a),o);if(s&&!l[i].test(h))throw new TypeError('Expected "'.concat(o.name,'" to match "').concat(o.pattern,'", but got "').concat(h,'"'));r+=o.prefix+h+o.suffix}}else r+=o}return r}}function i(t){return t.replace(/([.+*?=^!:${}()[\]|/\\])/g,"\\$1")}function o(t){return t&&t.sensitive?"":"i"}function n(e,r,n){return function(t,e,r){void 0===r&&(r={});for(var n=r.strict,a=void 0!==n&&n,s=r.start,l=void 0===s||s,c=r.end,d=void 0===c||c,p=r.encode,h=void 0===p?function(t){return t}:p,u=r.delimiter,m=void 0===u?"/#?":u,f=r.endsWith,v="[".concat(i(void 0===f?"":f),"]|$"),g="[".concat(i(m),"]"),y=l?"^":"",b=0,w=t;b<w.length;b++){var x=w[b];if("string"==typeof x)y+=i(h(x));else{var $=i(h(x.prefix)),k=i(h(x.suffix));if(x.pattern)if(e&&e.push(x),$||k)if("+"===x.modifier||"*"===x.modifier){var _="*"===x.modifier?"?":"";y+="(?:".concat($,"((?:").concat(x.pattern,")(?:").concat(k).concat($,"(?:").concat(x.pattern,"))*)").concat(k,")").concat(_)}else y+="(?:".concat($,"(").concat(x.pattern,")").concat(k,")").concat(x.modifier);else{if("+"===x.modifier||"*"===x.modifier)throw new TypeError('Can not repeat "'.concat(x.name,'" without a prefix and suffix'));y+="(".concat(x.pattern,")").concat(x.modifier)}else y+="(?:".concat($).concat(k,")").concat(x.modifier)}}if(d)a||(y+="".concat(g,"?")),y+=r.endsWith?"(?=".concat(v,")"):"$";else{var A=t[t.length-1],C="string"==typeof A?g.indexOf(A[A.length-1])>-1:void 0===A;a||(y+="(?:".concat(g,"(?=").concat(v,"))?")),C||(y+="(?=".concat(g,"|").concat(v,")"))}return new RegExp(y,o(r))}(t(e,n),r,n)}function a(t,e,r){return t instanceof RegExp?function(t,e){if(!e)return t;for(var r=/\((?:\?<(.*?)>)?(?!\?)/g,i=0,o=r.exec(t.source);o;)e.push({name:o[1]||i++,prefix:"",suffix:"",modifier:"",pattern:""}),o=r.exec(t.source);return t}(t,e):Array.isArray(t)?function(t,e,r){var i=t.map((function(t){return a(t,e,r).source}));return new RegExp("(?:".concat(i.join("|"),")"),o(r))}(t,e,r):n(t,e,r)}function s(t){return"object"==typeof t&&!!t}function l(t){return"function"==typeof t}function c(t){return"string"==typeof t}function d(t=[]){return Array.isArray(t)?t:[t]}function p(t){return`[Vaadin.Router] ${t}`}class h extends Error{code;context;constructor(t){super(p(`Page not found (${t.pathname})`)),this.context=t,this.code=404}}const u=Symbol("NotFoundResult");function m(t){return new h(t)}function f(t){return(Array.isArray(t)?t[0]:t)??""}function v(t){return f(t?.path)}const g=new Map;function y(t){try{return decodeURIComponent(t)}catch{return t}}g.set("|false",{keys:[],pattern:/(?:)/u});var b=function(t,e,r=!1,i=[],o){const n=`${t}|${String(r)}`,s=f(e);let l=g.get(n);if(!l){const e=[];l={keys:e,pattern:a(t,e,{end:r,strict:""===t})},g.set(n,l)}const c=l.pattern.exec(s);if(!c)return null;const d={...o};for(let t=1;t<c.length;t++){const e=l.keys[t-1],r=e.name,i=c[t];void 0===i&&Object.hasOwn(d,r)||("+"===e.modifier||"*"===e.modifier?d[r]=i?i.split(/[/?#]/u).map(y):[]:d[r]=i?y(i):i)}return{keys:[...i,...l.keys],params:d,path:c[0]}};var w=function t(e,r,i,o,n){let a,s,l=0,c=v(e);return c.startsWith("/")&&(i&&(c=c.substring(1)),i=!0),{next(d){if(e===d)return{done:!0,value:void 0};e.__children??=function(t){return Array.isArray(t)&&t.length>0?t:void 0}(e.children);const p=e.__children??[],h=!e.__children&&!e.children;if(!a&&(a=b(c,r,h,o,n),a))return{value:{keys:a.keys,params:a.params,path:a.path,route:e}};if(a&&p.length>0)for(;l<p.length;){if(!s){const o=p[l];o.parent=e;let n=a.path.length;n>0&&"/"===r.charAt(n)&&(n+=1),s=t(o,r.substring(n),i,a.keys,a.params)}const o=s.next(d);if(!o.done)return{done:!1,value:o.value};s=null,l+=1}return{done:!0,value:void 0}}}};function x(t){if(l(t.route.action))return t.route.action(t)}class $ extends Error{code;context;constructor(t,e){let r=`Path '${t.pathname}' is not properly resolved due to an error.`;const i=v(t.route);i&&(r+=` Resolution had failed on route: '${i}'`),super(r,e),this.code=e?.code,this.context=t}warn(){console.warn(this.message)}}class k{baseUrl;#t;errorHandler;resolveRoute;#e;constructor(t,{baseUrl:e="",context:r,errorHandler:i,resolveRoute:o=x}={}){if(Object(t)!==t)throw new TypeError("Invalid routes");this.baseUrl=e,this.errorHandler=i,this.resolveRoute=o,Array.isArray(t)?this.#e={__children:t,__synthetic:!0,action:()=>{},path:""}:this.#e={...t,parent:void 0},this.#t={...r,hash:"",next:async()=>u,params:{},pathname:"",resolver:this,route:this.#e,search:"",chain:[]}}get root(){return this.#e}get context(){return this.#t}get __effectiveBaseUrl(){return this.baseUrl?new URL(this.baseUrl,document.baseURI||document.URL).href.replace(/[^/]*$/u,""):""}getRoutes(){return[...this.#e.__children??[]]}removeRoutes(){this.#e.__children=[]}async resolve(t){const e=this,r={...this.#t,...c(t)?{pathname:t}:t,next:l},i=w(this.#e,this.__normalizePathname(r.pathname)??r.pathname,!!this.baseUrl),o=this.resolveRoute;let n=null,a=null,s=r;async function l(t=!1,c=n?.value?.route,d){const p=null===d?n?.value?.route:void 0;if(n=a??i.next(p),a=null,!t&&(n.done||!function(t,e){let r=t;for(;r;)if(r=r.parent,r===e)return!0;return!1}(n.value.route,c)))return a=n,u;if(n.done)throw m(r);s={...r,params:n.value.params,route:n.value.route,chain:s.chain?.slice()},function(t,e){const{path:r,route:i}=e;if(i&&!i.__synthetic){const e={path:r,route:i};if(i.parent&&t.chain)for(let e=t.chain.length-1;e>=0&&t.chain[e].route!==i.parent;e--)t.chain.pop();t.chain?.push(e)}}(s,n.value);const h=await o(s);return null!=h&&h!==u?(s.result=(f=h)&&"object"==typeof f&&"next"in f&&"params"in f&&"result"in f&&"route"in f?h.result:h,e.#t=s,s):await l(t,c,h);var f}try{return await l(!0,this.#e)}catch(t){const e=t instanceof h?t:new $(s,{code:500,cause:t});if(this.errorHandler)return s.result=this.errorHandler(e),s;throw t}}setRoutes(t){this.#e.__children=[...d(t)]}__normalizePathname(t){if(!this.baseUrl)return t;const e=this.__effectiveBaseUrl,r=t.startsWith("/")?new URL(e).origin+t:`./${t}`,i=new URL(r,e).href;return i.startsWith(e)?i.slice(e.length):void 0}addRoutes(t){return this.#e.__children=[...this.#e.__children??[],...d(t)],this.getRoutes()}}function _(t,e,r,i){const o=e.name??i?.(e);if(o&&(t.has(o)?t.get(o)?.push(e):t.set(o,[e])),Array.isArray(r))for(const o of r)o.parent=e,_(t,o,o.__children??o.children,i)}function A(t,e){const r=t.get(e);if(r){if(r.length>1)throw new Error(`Duplicate route with name "${e}". Try seting unique 'name' route properties.`);return r[0]}}var C=function(e,i={}){if(!(e instanceof k))throw new TypeError("An instance of Resolver is expected");const o=new Map,n=new Map;return(a,s)=>{let l=A(n,a);if(!l&&(n.clear(),_(n,e.root,e.root.__children,i.cacheKeyProvider),l=A(n,a),!l))throw new Error(`Route "${a}" not found`);let d=l.fullPath?o.get(l.fullPath):void 0;if(!d){let e=v(l),r=l.parent;for(;r;){const t=v(r);t&&(e=`${t.replace(/\/$/u,"")}/${e.replace(/^\//u,"")}`),r=r.parent}const i=t(e),n=Object.create(null);for(const t of i)c(t)||(n[t.name]=!0);d={keys:n,tokens:i},o.set(e,d),l.fullPath=e}let p=r(d.tokens,{encode:encodeURIComponent,...i})(s)||"/";if(i.stringifyQueryParams&&s){const t={};for(const[e,r]of Object.entries(s))!(e in d.keys)&&r&&(t[e]=r);const e=i.stringifyQueryParams(t);e&&(p+=e.startsWith("?")?e:`?${e}`)}return p}};const E=/\/\*[\*!]\s+vaadin-dev-mode:start([\s\S]*)vaadin-dev-mode:end\s+\*\*\//i,P=window.Vaadin&&window.Vaadin.Flow&&window.Vaadin.Flow.clients;function R(t,e){if("function"!=typeof t)return;const r=E.exec(t.toString());if(r)try{t=new Function(r[1])}catch(t){console.log("vaadin-development-mode-detector: uncommentAndRun() failed",t)}return t(e)}window.Vaadin=window.Vaadin||{};const O=function(t,e){if(window.Vaadin.developmentMode)return R(t,e)};function S(){
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

  vaadin-dev-mode:end **/}void 0===window.Vaadin.developmentMode&&(window.Vaadin.developmentMode=function(){try{return!!localStorage.getItem("vaadin.developmentmode.force")||["localhost","127.0.0.1"].indexOf(window.location.hostname)>=0&&(P?!(P&&Object.keys(P).map((t=>P[t])).filter((t=>t.productionMode)).length>0):!R((function(){return!0})))}catch(t){return!1}}());!function(t,e=(window.Vaadin??={})){e.registrations??=[],e.registrations.push({is:"@vaadin/router",version:"2.0.0"})}(),O(S);var z=async function(t,e){return t.classList.add(e),await new Promise((r=>{if((t=>{const e=getComputedStyle(t).getPropertyValue("animation-name");return e&&"none"!==e})(t)){const i=t.getBoundingClientRect(),o=`height: ${i.bottom-i.top}px; width: ${i.right-i.left}px`;t.setAttribute("style",`position: absolute; ${o}`),((t,e)=>{const r=()=>{t.removeEventListener("animationend",r),e()};t.addEventListener("animationend",r)})(t,(()=>{t.classList.remove(e),t.removeAttribute("style"),r()}))}else t.classList.remove(e),r()}))};function M(t){if(!t||!c(t.path))throw new Error(p('Expected route config to be an object with a "path" string property, or an array of such objects'));if(!(l(t.action)||Array.isArray(t.children)||l(t.children)||c(t.component)||c(t.redirect)))throw new Error(p(`Expected route config "${t.path}" to include either "component, redirect" or "action" function but none found.`));t.redirect&&["bundle","component"].forEach((e=>{e in t&&console.warn(p(`Route config "${String(t.path)}" has both "redirect" and "${e}" properties, and "redirect" will always override the latter. Did you mean to only use "${e}"?`))}))}function T(t){d(t).forEach((t=>M(t)))}function L(t,e){const r=e.__effectiveBaseUrl;return r?new URL(t.replace(/^\//u,""),r).pathname:t}function j(t){return t.map((t=>t.path)).reduce(((t,e)=>e.length?`${t.replace(/\/$/u,"")}/${e.replace(/^\//u,"")}`:t),"")}function B({chain:t=[],hash:r="",params:i={},pathname:o="",redirectFrom:n,resolver:a,search:s=""},l){const c=t.map((t=>t.route));return{baseUrl:a?.baseUrl??"",getUrl:(r={})=>a?L(e(function(t){return j(t.map((t=>t.route)))}(t))({...i,...r}),a):"",hash:r,params:i,pathname:o,redirectFrom:n,route:l??(Array.isArray(c)?c.at(-1):void 0)??null,routes:c,search:s,searchParams:new URLSearchParams(s)}}function I(t,e){const r={...t.params};return{redirect:{from:t.pathname,params:r,pathname:e}}}function N(t,e,...r){if("function"==typeof t)return t.apply(e,r)}function D(t,e,...r){return i=>i&&s(i)&&("cancel"in i||"redirect"in i)?i:N(e?.[t],e,...r)}function U(t,e){return!window.dispatchEvent(new CustomEvent(`vaadin-router-${t}`,{cancelable:"go"===t,detail:e}))}function H(t){if(t instanceof Element)return t.nodeName.toLowerCase()}function V(t){if(t.defaultPrevented)return;if(0!==t.button)return;if(t.shiftKey||t.ctrlKey||t.altKey||t.metaKey)return;let e=t.target;const r=t instanceof MouseEvent?t.composedPath():t.path??[];for(let t=0;t<r.length;t++){const i=r[t];if("nodeName"in i&&"a"===i.nodeName.toLowerCase()){e=i;break}}for(;e&&e instanceof Node&&"a"!==H(e);)e=e.parentNode;if(!e||"a"!==H(e))return;const i=e;if(i.target&&"_self"!==i.target.toLowerCase())return;if(i.hasAttribute("download"))return;if(i.hasAttribute("router-ignore"))return;if(i.pathname===window.location.pathname&&""!==i.hash)return;const o=i.origin||function(t){const{port:e,protocol:r}=t;return`${r}//${"http:"===r&&"80"===e||"https:"===r&&"443"===e?t.hostname:t.host}`}(i);if(o!==window.location.origin)return;const{hash:n,pathname:a,search:s}=i;U("go",{hash:n,pathname:a,search:s})&&t instanceof MouseEvent&&(t.preventDefault(),"click"===t.type&&window.scrollTo(0,0))}function K(t){if("vaadin-router-ignore"===t.state)return;const{hash:e,pathname:r,search:i}=window.location;U("go",{hash:e,pathname:r,search:i})}let F=[];const W={CLICK:{activate(){window.document.addEventListener("click",V)},inactivate(){window.document.removeEventListener("click",V)}},POPSTATE:{activate(){window.addEventListener("popstate",K)},inactivate(){window.removeEventListener("popstate",K)}}};function Y(t=[]){F.forEach((t=>t.inactivate())),t.forEach((t=>t.activate())),F=t}function G(){return{cancel:!0}}const q={__renderId:-1,params:{},route:{__synthetic:!0,children:[],path:"",action(){}},pathname:"",next:async()=>u};class J extends k{location=B({resolver:this});ready=Promise.resolve(this.location);#r=new WeakSet;#i=new WeakSet;#o=this.#n.bind(this);#a=0;#s;__previousContext;#l;#c=null;#d=null;constructor(t,e){const r=document.head.querySelector("base"),i=r?.getAttribute("href");super([],{baseUrl:i?new URL(i,document.URL).href.replace(/[^/]*$/u,""):void 0,...e,resolveRoute:async t=>await this.#p(t)}),Y(Object.values(W)),this.setOutlet(t),this.subscribe()}async#p(t){const{route:e}=t;if(l(e.children)){let r=await e.children(function({next:t,...e}){return e}(t));l(e.children)||({children:r}=e),function(t,e){if(!Array.isArray(t)&&!s(t))throw new Error(p(`Incorrect "children" value for the route ${String(e.path)}: expected array or object, but got ${String(t)}`));const r=d(t);r.forEach((t=>M(t))),e.__children=r}(r,e)}const r={component:t=>{const e=document.createElement(t);return this.#i.add(e),e},prevent:G,redirect:e=>I(t,e)};return await Promise.resolve().then((async()=>{if(this.#h(t))return await N(e.action,e,t,r)})).then((t=>null==t||"object"!=typeof t&&"symbol"!=typeof t||!(t instanceof HTMLElement||t===u||s(t)&&"redirect"in t)?c(e.redirect)?r.redirect(e.redirect):void 0:t)).then((t=>null!=t?t:c(e.component)?r.component(e.component):void 0))}setOutlet(t){t&&this.#u(t),this.#s=t}getOutlet(){return this.#s}async setRoutes(t,e=!1){return this.__previousContext=void 0,this.#l=void 0,T(t),super.setRoutes(t),e||this.#n(),await this.ready}addRoutes(t){return T(t),super.addRoutes(t)}async render(t,e=!1){this.#a+=1;const r=this.#a,i={...q,...c(t)?{hash:"",search:"",pathname:t}:t,__renderId:r};return this.ready=this.#m(i,e),await this.ready}async#m(t,e){const{__renderId:r}=t;try{const i=await this.resolve(t),o=await this.#f(i);if(!this.#h(o))return this.location;const n=this.__previousContext;if(o===n)return this.#v(n,!0),this.location;if(this.location=B(o),e&&this.#v(o,1===r),U("location-changed",{router:this,location:this.location}),o.__skipAttach)return this.#g(o,n),this.__previousContext=o,this.location;this.#y(o,n);const a=this.#b(o);if(this.#w(o),this.#x(o,n),await a,this.#h(o))return this.#$(),this.__previousContext=o,this.location}catch(i){if(r===this.#a){e&&this.#v(this.context);for(const t of this.#s?.children??[])t.remove();throw this.location=B(Object.assign(t,{resolver:this})),U("error",{router:this,error:i,...t}),i}}return this.location}async#f(t,e=t){const r=await this.#k(e),i=r!==e?r:t,o=L(j(r.chain??[]),this)===r.pathname,n=async(t,e=t.route,r)=>{const i=await t.next(!1,e,r);return null===i||i===u?o?t:null!=e.parent?await n(t,e.parent,i):i:i},a=await n(r);if(null==a||a===u)throw m(i);return a!==r?await this.#f(i,a):await this.#_(r)}async#k(t){const{result:e}=t;if(e instanceof HTMLElement)return function(t,e){if(e.location=B(t),t.chain){const r=t.chain.map((t=>t.route)).indexOf(t.route);t.chain[r].element=e}}(t,e),t;if(e&&"redirect"in e){const r=await this.#A(e.redirect,t.__redirectCount,t.__renderId);return await this.#k(r)}throw e instanceof Error?e:new Error(p(`Invalid route resolution result for path "${t.pathname}". Expected redirect object or HTML element, but got: "${function(t){if("object"!=typeof t)return String(t);const[e="Unknown"]=/ (.*)\]$/u.exec(String(t))??[];return"Object"===e||"Array"===e?`${e} ${JSON.stringify(t)}`:e}(e)}". Double check the action return value for the route.`))}async#_(t){return await this.#C(t).then((async e=>e===this.__previousContext||e===t?e:await this.#f(e)))}async#C(t){const e=this.__previousContext??{},r=e.chain??[],i=t.chain??[];let o=Promise.resolve(void 0);const n=e=>I(t,e);if(t.__divergedChainIndex=0,t.__skipAttach=!1,r.length){for(let e=0;e<Math.min(r.length,i.length)&&(r[e].route===i[e].route&&(r[e].path===i[e].path||r[e].element===i[e].element)&&this.#E(r[e].element,i[e].element));t.__divergedChainIndex++,e++);if(t.__skipAttach=i.length===r.length&&t.__divergedChainIndex===i.length&&this.#E(t.result,e.result),t.__skipAttach){for(let e=i.length-1;e>=0;e--)o=this.#P(o,t,{prevent:G},r[e]);for(let e=0;e<i.length;e++)o=this.#R(o,t,{prevent:G,redirect:n},i[e]),r[e].element.location=B(t,r[e].route)}else for(let e=r.length-1;e>=t.__divergedChainIndex;e--)o=this.#P(o,t,{prevent:G},r[e])}if(!t.__skipAttach)for(let e=0;e<i.length;e++)e<t.__divergedChainIndex?e<r.length&&r[e].element&&(r[e].element.location=B(t,r[e].route)):(o=this.#R(o,t,{prevent:G,redirect:n},i[e]),i[e].element&&(i[e].element.location=B(t,i[e].route)));return await o.then((async e=>{if(e&&s(e)){if("cancel"in e&&this.__previousContext)return this.__previousContext.__renderId=t.__renderId,this.__previousContext;if("redirect"in e)return await this.#A(e.redirect,t.__redirectCount,t.__renderId)}return t}))}async#P(t,e,r,i){const o=B(e);let n=await t;if(this.#h(e)){n=D("onBeforeLeave",i.element,o,r,this)(n)}if(!s(n)||!("redirect"in n))return n}async#R(t,e,r,i){const o=B(e,i.route),n=await t;if(this.#h(e)){return D("onBeforeEnter",i.element,o,r,this)(n)}}#E(t,e){return t instanceof Element&&e instanceof Element&&(this.#i.has(t)&&this.#i.has(e)?t.localName===e.localName:t===e)}#h(t){return t.__renderId===this.#a}async#A(t,e=0,r=0){if(e>256)throw new Error(p(`Too many redirects when rendering ${t.from}`));return await this.resolve({...q,pathname:this.urlForPath(t.pathname,t.params),redirectFrom:t.from,__redirectCount:e+1,__renderId:r})}#u(t=this.#s){if(!(t instanceof Element||t instanceof DocumentFragment))throw new TypeError(p(`Expected router outlet to be a valid DOM Element | DocumentFragment (but got ${t})`))}#v({pathname:t,search:e="",hash:r=""},i){if(window.location.pathname!==t||window.location.search!==e||window.location.hash!==r){const o=i?"replaceState":"pushState";window.history[o](null,document.title,t+e+r),window.dispatchEvent(new PopStateEvent("popstate",{state:"vaadin-router-ignore"}))}}#g(t,e){let r=this.#s;for(let i=0;i<(t.__divergedChainIndex??0);i++){const o=e?.chain?.[i].element;if(o){if(o.parentNode!==r)break;t.chain[i].element=o,r=o}}return r}#y(t,e){this.#u(),this.#O();const r=this.#g(t,e);this.#c=[],this.#d=Array.from(r?.children??[]).filter((e=>this.#r.has(e)&&e!==t.result));let i=r;for(let e=t.__divergedChainIndex??0;e<(t.chain?.length??0);e++){const o=t.chain[e].element;o&&(i?.appendChild(o),this.#r.add(o),i===r&&this.#c.push(o),i=o)}}#$(){if(this.#d)for(const t of this.#d)t.remove();this.#d=null,this.#c=null}#O(){if(this.#d&&this.#c){for(const t of this.#c)t.remove();this.#d=null,this.#c=null}}#x(t,e){if(e?.chain&&null!=t.__divergedChainIndex)for(let r=e.chain.length-1;r>=t.__divergedChainIndex&&this.#h(t);r--){const i=e.chain[r].element;if(i)try{const e=B(t);N(i.onAfterLeave,i,e,{},this)}finally{if(this.#d?.includes(i))for(const t of i.children)t.remove()}}}#w(t){if(t.chain&&null!=t.__divergedChainIndex)for(let e=t.__divergedChainIndex;e<t.chain.length&&this.#h(t);e++){const r=t.chain[e].element;if(r){const i=B(t,t.chain[e].route);N(r.onAfterEnter,r,i,{},this)}}}async#b(t){const e=this.#d?.[0],r=this.#c?.[0],i=[],{chain:o=[]}=t;let n;for(let t=o.length-1;t>=0;t--)if(o[t].route.animate){n=o[t].route.animate;break}if(e&&r&&n){const t=s(n)&&n.leave?n.leave:"leaving",o=s(n)&&n.enter?n.enter:"entering";i.push(z(e,t)),i.push(z(r,o))}return await Promise.all(i),t}subscribe(){window.addEventListener("vaadin-router-go",this.#o)}unsubscribe(){window.removeEventListener("vaadin-router-go",this.#o)}#n(t){const{pathname:e,search:r,hash:i}=t instanceof CustomEvent?t.detail:window.location;c(this.__normalizePathname(e))&&(t?.preventDefault&&t.preventDefault(),this.render({pathname:e,search:r,hash:i},!0))}static setTriggers(...t){Y(t)}urlForName(t,e){return this.#l||(this.#l=C(this,{cacheKeyProvider:t=>"component"in t&&"string"==typeof t.component?t.component:void 0})),L(this.#l(t,e??void 0),this)}urlForPath(t,r){return L(e(t)(r??void 0),this)}static go(t){const{pathname:e,search:r,hash:i}=c(t)?new URL(t,"http://a"):t;return U("go",{pathname:e,search:r,hash:i})}}
/**
 * @license
 * Copyright 2019 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */const Q=globalThis,Z=Q.ShadowRoot&&(void 0===Q.ShadyCSS||Q.ShadyCSS.nativeShadow)&&"adoptedStyleSheets"in Document.prototype&&"replace"in CSSStyleSheet.prototype,X=Symbol(),tt=new WeakMap;let et=class{constructor(t,e,r){if(this._$cssResult$=!0,r!==X)throw Error("CSSResult is not constructable. Use `unsafeCSS` or `css` instead.");this.cssText=t,this.t=e}get styleSheet(){let t=this.o;const e=this.t;if(Z&&void 0===t){const r=void 0!==e&&1===e.length;r&&(t=tt.get(e)),void 0===t&&((this.o=t=new CSSStyleSheet).replaceSync(this.cssText),r&&tt.set(e,t))}return t}toString(){return this.cssText}};const rt=(t,...e)=>{const r=1===t.length?t[0]:e.reduce(((e,r,i)=>e+(t=>{if(!0===t._$cssResult$)return t.cssText;if("number"==typeof t)return t;throw Error("Value passed to 'css' function must be a 'css' function result: "+t+". Use 'unsafeCSS' to pass non-literal values, but take care to ensure page security.")})(r)+t[i+1]),t[0]);return new et(r,t,X)},it=Z?t=>t:t=>t instanceof CSSStyleSheet?(t=>{let e="";for(const r of t.cssRules)e+=r.cssText;return(t=>new et("string"==typeof t?t:t+"",void 0,X))(e)})(t):t
/**
 * @license
 * Copyright 2017 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */,{is:ot,defineProperty:nt,getOwnPropertyDescriptor:at,getOwnPropertyNames:st,getOwnPropertySymbols:lt,getPrototypeOf:ct}=Object,dt=globalThis,pt=dt.trustedTypes,ht=pt?pt.emptyScript:"",ut=dt.reactiveElementPolyfillSupport,mt=(t,e)=>t,ft={toAttribute(t,e){switch(e){case Boolean:t=t?ht:null;break;case Object:case Array:t=null==t?t:JSON.stringify(t)}return t},fromAttribute(t,e){let r=t;switch(e){case Boolean:r=null!==t;break;case Number:r=null===t?null:Number(t);break;case Object:case Array:try{r=JSON.parse(t)}catch(t){r=null}}return r}},vt=(t,e)=>!ot(t,e),gt={attribute:!0,type:String,converter:ft,reflect:!1,hasChanged:vt};Symbol.metadata??=Symbol("metadata"),dt.litPropertyMetadata??=new WeakMap;class yt extends HTMLElement{static addInitializer(t){this._$Ei(),(this.l??=[]).push(t)}static get observedAttributes(){return this.finalize(),this._$Eh&&[...this._$Eh.keys()]}static createProperty(t,e=gt){if(e.state&&(e.attribute=!1),this._$Ei(),this.elementProperties.set(t,e),!e.noAccessor){const r=Symbol(),i=this.getPropertyDescriptor(t,r,e);void 0!==i&&nt(this.prototype,t,i)}}static getPropertyDescriptor(t,e,r){const{get:i,set:o}=at(this.prototype,t)??{get(){return this[e]},set(t){this[e]=t}};return{get(){return i?.call(this)},set(e){const n=i?.call(this);o.call(this,e),this.requestUpdate(t,n,r)},configurable:!0,enumerable:!0}}static getPropertyOptions(t){return this.elementProperties.get(t)??gt}static _$Ei(){if(this.hasOwnProperty(mt("elementProperties")))return;const t=ct(this);t.finalize(),void 0!==t.l&&(this.l=[...t.l]),this.elementProperties=new Map(t.elementProperties)}static finalize(){if(this.hasOwnProperty(mt("finalized")))return;if(this.finalized=!0,this._$Ei(),this.hasOwnProperty(mt("properties"))){const t=this.properties,e=[...st(t),...lt(t)];for(const r of e)this.createProperty(r,t[r])}const t=this[Symbol.metadata];if(null!==t){const e=litPropertyMetadata.get(t);if(void 0!==e)for(const[t,r]of e)this.elementProperties.set(t,r)}this._$Eh=new Map;for(const[t,e]of this.elementProperties){const r=this._$Eu(t,e);void 0!==r&&this._$Eh.set(r,t)}this.elementStyles=this.finalizeStyles(this.styles)}static finalizeStyles(t){const e=[];if(Array.isArray(t)){const r=new Set(t.flat(1/0).reverse());for(const t of r)e.unshift(it(t))}else void 0!==t&&e.push(it(t));return e}static _$Eu(t,e){const r=e.attribute;return!1===r?void 0:"string"==typeof r?r:"string"==typeof t?t.toLowerCase():void 0}constructor(){super(),this._$Ep=void 0,this.isUpdatePending=!1,this.hasUpdated=!1,this._$Em=null,this._$Ev()}_$Ev(){this._$ES=new Promise((t=>this.enableUpdating=t)),this._$AL=new Map,this._$E_(),this.requestUpdate(),this.constructor.l?.forEach((t=>t(this)))}addController(t){(this._$EO??=new Set).add(t),void 0!==this.renderRoot&&this.isConnected&&t.hostConnected?.()}removeController(t){this._$EO?.delete(t)}_$E_(){const t=new Map,e=this.constructor.elementProperties;for(const r of e.keys())this.hasOwnProperty(r)&&(t.set(r,this[r]),delete this[r]);t.size>0&&(this._$Ep=t)}createRenderRoot(){const t=this.shadowRoot??this.attachShadow(this.constructor.shadowRootOptions);return((t,e)=>{if(Z)t.adoptedStyleSheets=e.map((t=>t instanceof CSSStyleSheet?t:t.styleSheet));else for(const r of e){const e=document.createElement("style"),i=Q.litNonce;void 0!==i&&e.setAttribute("nonce",i),e.textContent=r.cssText,t.appendChild(e)}})(t,this.constructor.elementStyles),t}connectedCallback(){this.renderRoot??=this.createRenderRoot(),this.enableUpdating(!0),this._$EO?.forEach((t=>t.hostConnected?.()))}enableUpdating(t){}disconnectedCallback(){this._$EO?.forEach((t=>t.hostDisconnected?.()))}attributeChangedCallback(t,e,r){this._$AK(t,r)}_$EC(t,e){const r=this.constructor.elementProperties.get(t),i=this.constructor._$Eu(t,r);if(void 0!==i&&!0===r.reflect){const o=(void 0!==r.converter?.toAttribute?r.converter:ft).toAttribute(e,r.type);this._$Em=t,null==o?this.removeAttribute(i):this.setAttribute(i,o),this._$Em=null}}_$AK(t,e){const r=this.constructor,i=r._$Eh.get(t);if(void 0!==i&&this._$Em!==i){const t=r.getPropertyOptions(i),o="function"==typeof t.converter?{fromAttribute:t.converter}:void 0!==t.converter?.fromAttribute?t.converter:ft;this._$Em=i,this[i]=o.fromAttribute(e,t.type),this._$Em=null}}requestUpdate(t,e,r){if(void 0!==t){if(r??=this.constructor.getPropertyOptions(t),!(r.hasChanged??vt)(this[t],e))return;this.P(t,e,r)}!1===this.isUpdatePending&&(this._$ES=this._$ET())}P(t,e,r){this._$AL.has(t)||this._$AL.set(t,e),!0===r.reflect&&this._$Em!==t&&(this._$Ej??=new Set).add(t)}async _$ET(){this.isUpdatePending=!0;try{await this._$ES}catch(t){Promise.reject(t)}const t=this.scheduleUpdate();return null!=t&&await t,!this.isUpdatePending}scheduleUpdate(){return this.performUpdate()}performUpdate(){if(!this.isUpdatePending)return;if(!this.hasUpdated){if(this.renderRoot??=this.createRenderRoot(),this._$Ep){for(const[t,e]of this._$Ep)this[t]=e;this._$Ep=void 0}const t=this.constructor.elementProperties;if(t.size>0)for(const[e,r]of t)!0!==r.wrapped||this._$AL.has(e)||void 0===this[e]||this.P(e,this[e],r)}let t=!1;const e=this._$AL;try{t=this.shouldUpdate(e),t?(this.willUpdate(e),this._$EO?.forEach((t=>t.hostUpdate?.())),this.update(e)):this._$EU()}catch(e){throw t=!1,this._$EU(),e}t&&this._$AE(e)}willUpdate(t){}_$AE(t){this._$EO?.forEach((t=>t.hostUpdated?.())),this.hasUpdated||(this.hasUpdated=!0,this.firstUpdated(t)),this.updated(t)}_$EU(){this._$AL=new Map,this.isUpdatePending=!1}get updateComplete(){return this.getUpdateComplete()}getUpdateComplete(){return this._$ES}shouldUpdate(t){return!0}update(t){this._$Ej&&=this._$Ej.forEach((t=>this._$EC(t,this[t]))),this._$EU()}updated(t){}firstUpdated(t){}}yt.elementStyles=[],yt.shadowRootOptions={mode:"open"},yt[mt("elementProperties")]=new Map,yt[mt("finalized")]=new Map,ut?.({ReactiveElement:yt}),(dt.reactiveElementVersions??=[]).push("2.0.4");
/**
 * @license
 * Copyright 2017 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */
const bt=globalThis,wt=bt.trustedTypes,xt=wt?wt.createPolicy("lit-html",{createHTML:t=>t}):void 0,$t="$lit$",kt=`lit$${Math.random().toFixed(9).slice(2)}$`,_t="?"+kt,At=`<${_t}>`,Ct=document,Et=()=>Ct.createComment(""),Pt=t=>null===t||"object"!=typeof t&&"function"!=typeof t,Rt=Array.isArray,Ot="[ \t\n\f\r]",St=/<(?:(!--|\/[^a-zA-Z])|(\/?[a-zA-Z][^>\s]*)|(\/?$))/g,zt=/-->/g,Mt=/>/g,Tt=RegExp(`>|${Ot}(?:([^\\s"'>=/]+)(${Ot}*=${Ot}*(?:[^ \t\n\f\r"'\`<>=]|("|')|))|$)`,"g"),Lt=/'/g,jt=/"/g,Bt=/^(?:script|style|textarea|title)$/i,It=(t=>(e,...r)=>({_$litType$:t,strings:e,values:r}))(1),Nt=Symbol.for("lit-noChange"),Dt=Symbol.for("lit-nothing"),Ut=new WeakMap,Ht=Ct.createTreeWalker(Ct,129);function Vt(t,e){if(!Rt(t)||!t.hasOwnProperty("raw"))throw Error("invalid template strings array");return void 0!==xt?xt.createHTML(e):e}class Kt{constructor({strings:t,_$litType$:e},r){let i;this.parts=[];let o=0,n=0;const a=t.length-1,s=this.parts,[l,c]=((t,e)=>{const r=t.length-1,i=[];let o,n=2===e?"<svg>":3===e?"<math>":"",a=St;for(let e=0;e<r;e++){const r=t[e];let s,l,c=-1,d=0;for(;d<r.length&&(a.lastIndex=d,l=a.exec(r),null!==l);)d=a.lastIndex,a===St?"!--"===l[1]?a=zt:void 0!==l[1]?a=Mt:void 0!==l[2]?(Bt.test(l[2])&&(o=RegExp("</"+l[2],"g")),a=Tt):void 0!==l[3]&&(a=Tt):a===Tt?">"===l[0]?(a=o??St,c=-1):void 0===l[1]?c=-2:(c=a.lastIndex-l[2].length,s=l[1],a=void 0===l[3]?Tt:'"'===l[3]?jt:Lt):a===jt||a===Lt?a=Tt:a===zt||a===Mt?a=St:(a=Tt,o=void 0);const p=a===Tt&&t[e+1].startsWith("/>")?" ":"";n+=a===St?r+At:c>=0?(i.push(s),r.slice(0,c)+$t+r.slice(c)+kt+p):r+kt+(-2===c?e:p)}return[Vt(t,n+(t[r]||"<?>")+(2===e?"</svg>":3===e?"</math>":"")),i]})(t,e);if(this.el=Kt.createElement(l,r),Ht.currentNode=this.el.content,2===e||3===e){const t=this.el.content.firstChild;t.replaceWith(...t.childNodes)}for(;null!==(i=Ht.nextNode())&&s.length<a;){if(1===i.nodeType){if(i.hasAttributes())for(const t of i.getAttributeNames())if(t.endsWith($t)){const e=c[n++],r=i.getAttribute(t).split(kt),a=/([.?@])?(.*)/.exec(e);s.push({type:1,index:o,name:a[2],strings:r,ctor:"."===a[1]?qt:"?"===a[1]?Jt:"@"===a[1]?Qt:Gt}),i.removeAttribute(t)}else t.startsWith(kt)&&(s.push({type:6,index:o}),i.removeAttribute(t));if(Bt.test(i.tagName)){const t=i.textContent.split(kt),e=t.length-1;if(e>0){i.textContent=wt?wt.emptyScript:"";for(let r=0;r<e;r++)i.append(t[r],Et()),Ht.nextNode(),s.push({type:2,index:++o});i.append(t[e],Et())}}}else if(8===i.nodeType)if(i.data===_t)s.push({type:2,index:o});else{let t=-1;for(;-1!==(t=i.data.indexOf(kt,t+1));)s.push({type:7,index:o}),t+=kt.length-1}o++}}static createElement(t,e){const r=Ct.createElement("template");return r.innerHTML=t,r}}function Ft(t,e,r=t,i){if(e===Nt)return e;let o=void 0!==i?r._$Co?.[i]:r._$Cl;const n=Pt(e)?void 0:e._$litDirective$;return o?.constructor!==n&&(o?._$AO?.(!1),void 0===n?o=void 0:(o=new n(t),o._$AT(t,r,i)),void 0!==i?(r._$Co??=[])[i]=o:r._$Cl=o),void 0!==o&&(e=Ft(t,o._$AS(t,e.values),o,i)),e}let Wt=class{constructor(t,e){this._$AV=[],this._$AN=void 0,this._$AD=t,this._$AM=e}get parentNode(){return this._$AM.parentNode}get _$AU(){return this._$AM._$AU}u(t){const{el:{content:e},parts:r}=this._$AD,i=(t?.creationScope??Ct).importNode(e,!0);Ht.currentNode=i;let o=Ht.nextNode(),n=0,a=0,s=r[0];for(;void 0!==s;){if(n===s.index){let e;2===s.type?e=new Yt(o,o.nextSibling,this,t):1===s.type?e=new s.ctor(o,s.name,s.strings,this,t):6===s.type&&(e=new Zt(o,this,t)),this._$AV.push(e),s=r[++a]}n!==s?.index&&(o=Ht.nextNode(),n++)}return Ht.currentNode=Ct,i}p(t){let e=0;for(const r of this._$AV)void 0!==r&&(void 0!==r.strings?(r._$AI(t,r,e),e+=r.strings.length-2):r._$AI(t[e])),e++}};class Yt{get _$AU(){return this._$AM?._$AU??this._$Cv}constructor(t,e,r,i){this.type=2,this._$AH=Dt,this._$AN=void 0,this._$AA=t,this._$AB=e,this._$AM=r,this.options=i,this._$Cv=i?.isConnected??!0}get parentNode(){let t=this._$AA.parentNode;const e=this._$AM;return void 0!==e&&11===t?.nodeType&&(t=e.parentNode),t}get startNode(){return this._$AA}get endNode(){return this._$AB}_$AI(t,e=this){t=Ft(this,t,e),Pt(t)?t===Dt||null==t||""===t?(this._$AH!==Dt&&this._$AR(),this._$AH=Dt):t!==this._$AH&&t!==Nt&&this._(t):void 0!==t._$litType$?this.$(t):void 0!==t.nodeType?this.T(t):(t=>Rt(t)||"function"==typeof t?.[Symbol.iterator])(t)?this.k(t):this._(t)}O(t){return this._$AA.parentNode.insertBefore(t,this._$AB)}T(t){this._$AH!==t&&(this._$AR(),this._$AH=this.O(t))}_(t){this._$AH!==Dt&&Pt(this._$AH)?this._$AA.nextSibling.data=t:this.T(Ct.createTextNode(t)),this._$AH=t}$(t){const{values:e,_$litType$:r}=t,i="number"==typeof r?this._$AC(t):(void 0===r.el&&(r.el=Kt.createElement(Vt(r.h,r.h[0]),this.options)),r);if(this._$AH?._$AD===i)this._$AH.p(e);else{const t=new Wt(i,this),r=t.u(this.options);t.p(e),this.T(r),this._$AH=t}}_$AC(t){let e=Ut.get(t.strings);return void 0===e&&Ut.set(t.strings,e=new Kt(t)),e}k(t){Rt(this._$AH)||(this._$AH=[],this._$AR());const e=this._$AH;let r,i=0;for(const o of t)i===e.length?e.push(r=new Yt(this.O(Et()),this.O(Et()),this,this.options)):r=e[i],r._$AI(o),i++;i<e.length&&(this._$AR(r&&r._$AB.nextSibling,i),e.length=i)}_$AR(t=this._$AA.nextSibling,e){for(this._$AP?.(!1,!0,e);t&&t!==this._$AB;){const e=t.nextSibling;t.remove(),t=e}}setConnected(t){void 0===this._$AM&&(this._$Cv=t,this._$AP?.(t))}}class Gt{get tagName(){return this.element.tagName}get _$AU(){return this._$AM._$AU}constructor(t,e,r,i,o){this.type=1,this._$AH=Dt,this._$AN=void 0,this.element=t,this.name=e,this._$AM=i,this.options=o,r.length>2||""!==r[0]||""!==r[1]?(this._$AH=Array(r.length-1).fill(new String),this.strings=r):this._$AH=Dt}_$AI(t,e=this,r,i){const o=this.strings;let n=!1;if(void 0===o)t=Ft(this,t,e,0),n=!Pt(t)||t!==this._$AH&&t!==Nt,n&&(this._$AH=t);else{const i=t;let a,s;for(t=o[0],a=0;a<o.length-1;a++)s=Ft(this,i[r+a],e,a),s===Nt&&(s=this._$AH[a]),n||=!Pt(s)||s!==this._$AH[a],s===Dt?t=Dt:t!==Dt&&(t+=(s??"")+o[a+1]),this._$AH[a]=s}n&&!i&&this.j(t)}j(t){t===Dt?this.element.removeAttribute(this.name):this.element.setAttribute(this.name,t??"")}}class qt extends Gt{constructor(){super(...arguments),this.type=3}j(t){this.element[this.name]=t===Dt?void 0:t}}class Jt extends Gt{constructor(){super(...arguments),this.type=4}j(t){this.element.toggleAttribute(this.name,!!t&&t!==Dt)}}class Qt extends Gt{constructor(t,e,r,i,o){super(t,e,r,i,o),this.type=5}_$AI(t,e=this){if((t=Ft(this,t,e,0)??Dt)===Nt)return;const r=this._$AH,i=t===Dt&&r!==Dt||t.capture!==r.capture||t.once!==r.once||t.passive!==r.passive,o=t!==Dt&&(r===Dt||i);i&&this.element.removeEventListener(this.name,this,r),o&&this.element.addEventListener(this.name,this,t),this._$AH=t}handleEvent(t){"function"==typeof this._$AH?this._$AH.call(this.options?.host??this.element,t):this._$AH.handleEvent(t)}}class Zt{constructor(t,e,r){this.element=t,this.type=6,this._$AN=void 0,this._$AM=e,this.options=r}get _$AU(){return this._$AM._$AU}_$AI(t){Ft(this,t)}}const Xt={I:Yt},te=bt.litHtmlPolyfillSupport;te?.(Kt,Yt),(bt.litHtmlVersions??=[]).push("3.2.1");
/**
 * @license
 * Copyright 2017 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */
let ee=class extends yt{constructor(){super(...arguments),this.renderOptions={host:this},this._$Do=void 0}createRenderRoot(){const t=super.createRenderRoot();return this.renderOptions.renderBefore??=t.firstChild,t}update(t){const e=this.render();this.hasUpdated||(this.renderOptions.isConnected=this.isConnected),super.update(t),this._$Do=((t,e,r)=>{const i=r?.renderBefore??e;let o=i._$litPart$;if(void 0===o){const t=r?.renderBefore??null;i._$litPart$=o=new Yt(e.insertBefore(Et(),t),t,void 0,r??{})}return o._$AI(t),o})(e,this.renderRoot,this.renderOptions)}connectedCallback(){super.connectedCallback(),this._$Do?.setConnected(!0)}disconnectedCallback(){super.disconnectedCallback(),this._$Do?.setConnected(!1)}render(){return Nt}};ee._$litElement$=!0,ee.finalized=!0,globalThis.litElementHydrateSupport?.({LitElement:ee});const re=globalThis.litElementPolyfillSupport;re?.({LitElement:ee}),(globalThis.litElementVersions??=[]).push("4.1.1");
/**
 * @license
 * Copyright 2017 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */
const ie=t=>(e,r)=>{void 0!==r?r.addInitializer((()=>{customElements.define(t,e)})):customElements.define(t,e)}
/**
 * @license
 * Copyright 2017 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */,oe={attribute:!0,type:String,converter:ft,reflect:!1,hasChanged:vt},ne=(t=oe,e,r)=>{const{kind:i,metadata:o}=r;let n=globalThis.litPropertyMetadata.get(o);if(void 0===n&&globalThis.litPropertyMetadata.set(o,n=new Map),n.set(r.name,t),"accessor"===i){const{name:i}=r;return{set(r){const o=e.get.call(this);e.set.call(this,r),this.requestUpdate(i,o,t)},init(e){return void 0!==e&&this.P(i,void 0,t),e}}}if("setter"===i){const{name:i}=r;return function(r){const o=this[i];e.call(this,r),this.requestUpdate(i,o,t)}}throw Error("Unsupported decorator location: "+i)};function ae(t){return(e,r)=>"object"==typeof r?ne(t,e,r):((t,e,r)=>{const i=e.hasOwnProperty(r);return e.constructor.createProperty(r,i?{...t,wrapped:!0}:t),i?Object.getOwnPropertyDescriptor(e,r):void 0})(t,e,r)
/**
 * @license
 * Copyright 2017 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */}function se(t){return ae({...t,state:!0,attribute:!1})}var le=rt`
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
 */;const{I:ce}=Xt,de=()=>document.createComment(""),pe=(t,e,r)=>{const i=t._$AA.parentNode,o=void 0===e?t._$AB:e._$AA;if(void 0===r){const e=i.insertBefore(de(),o),n=i.insertBefore(de(),o);r=new ce(e,n,t,t.options)}else{const e=r._$AB.nextSibling,n=r._$AM,a=n!==t;if(a){let e;r._$AQ?.(t),r._$AM=t,void 0!==r._$AP&&(e=t._$AU)!==n._$AU&&r._$AP(e)}if(e!==o||a){let t=r._$AA;for(;t!==e;){const e=t.nextSibling;i.insertBefore(t,o),t=e}}}return r},he=(t,e,r=t)=>(t._$AI(e,r),t),ue={},me=t=>{t._$AP?.(!1,!0);let e=t._$AA;const r=t._$AB.nextSibling;for(;e!==r;){const t=e.nextSibling;e.remove(),e=t}},fe=1,ve=2,ge=t=>(...e)=>({_$litDirective$:t,values:e});
/**
 * @license
 * Copyright 2017 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */class ye{constructor(t){}get _$AU(){return this._$AM._$AU}_$AT(t,e,r){this._$Ct=t,this._$AM=e,this._$Ci=r}_$AS(t,e){return this.update(t,e)}update(t,e){return this.render(...e)}}
/**
 * @license
 * Copyright 2017 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */const be=(t,e)=>{const r=t._$AN;if(void 0===r)return!1;for(const t of r)t._$AO?.(e,!1),be(t,e);return!0},we=t=>{let e,r;do{if(void 0===(e=t._$AM))break;r=e._$AN,r.delete(t),t=e}while(0===r?.size)},xe=t=>{for(let e;e=t._$AM;t=e){let r=e._$AN;if(void 0===r)e._$AN=r=new Set;else if(r.has(t))break;r.add(t),_e(e)}};function $e(t){void 0!==this._$AN?(we(this),this._$AM=t,xe(this)):this._$AM=t}function ke(t,e=!1,r=0){const i=this._$AH,o=this._$AN;if(void 0!==o&&0!==o.size)if(e)if(Array.isArray(i))for(let t=r;t<i.length;t++)be(i[t],!1),we(i[t]);else null!=i&&(be(i,!1),we(i));else be(this,t)}const _e=t=>{t.type==ve&&(t._$AP??=ke,t._$AQ??=$e)};class Ae extends ye{constructor(){super(...arguments),this._$AN=void 0}_$AT(t,e,r){super._$AT(t,e,r),xe(this),this.isConnected=t._$AU}_$AO(t,e=!0){t!==this.isConnected&&(this.isConnected=t,t?this.reconnected?.():this.disconnected?.()),e&&(be(this,t),we(this))}setValue(t){if((t=>void 0===t.strings)(this._$Ct))this._$Ct._$AI(t,this);else{const e=[...this._$Ct._$AH];e[this._$Ci]=t,this._$Ct._$AI(e,this,0)}}disconnected(){}reconnected(){}}
/**
 * @license
 * Copyright 2020 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */const Ce=()=>new Ee;class Ee{}const Pe=new WeakMap,Re=ge(class extends Ae{render(t){return Dt}update(t,[e]){const r=e!==this.Y;return r&&void 0!==this.Y&&this.rt(void 0),(r||this.lt!==this.ct)&&(this.Y=e,this.ht=t.options?.host,this.rt(this.ct=t.element)),Dt}rt(t){if(this.isConnected||(t=void 0),"function"==typeof this.Y){const e=this.ht??globalThis;let r=Pe.get(e);void 0===r&&(r=new WeakMap,Pe.set(e,r)),void 0!==r.get(this.Y)&&this.Y.call(this.ht,void 0),r.set(this.Y,t),void 0!==t&&this.Y.call(this.ht,t)}else this.Y.value=t}get lt(){return"function"==typeof this.Y?Pe.get(this.ht??globalThis)?.get(this.Y):this.Y?.value}disconnected(){this.lt===this.ct&&this.rt(void 0)}reconnected(){this.rt(this.ct)}}),Oe=(t,e,r)=>{const i=new Map;for(let o=e;o<=r;o++)i.set(t[o],o);return i},Se=ge(class extends ye{constructor(t){if(super(t),t.type!==ve)throw Error("repeat() can only be used in text expressions")}dt(t,e,r){let i;void 0===r?r=e:void 0!==e&&(i=e);const o=[],n=[];let a=0;for(const e of t)o[a]=i?i(e,a):a,n[a]=r(e,a),a++;return{values:n,keys:o}}render(t,e,r){return this.dt(t,e,r).values}update(t,[e,r,i]){const o=(t=>t._$AH)(t),{values:n,keys:a}=this.dt(e,r,i);if(!Array.isArray(o))return this.ut=a,n;const s=this.ut??=[],l=[];let c,d,p=0,h=o.length-1,u=0,m=n.length-1;for(;p<=h&&u<=m;)if(null===o[p])p++;else if(null===o[h])h--;else if(s[p]===a[u])l[u]=he(o[p],n[u]),p++,u++;else if(s[h]===a[m])l[m]=he(o[h],n[m]),h--,m--;else if(s[p]===a[m])l[m]=he(o[p],n[m]),pe(t,l[m+1],o[p]),p++,m--;else if(s[h]===a[u])l[u]=he(o[h],n[u]),pe(t,o[p],o[h]),h--,u++;else if(void 0===c&&(c=Oe(a,u,m),d=Oe(s,p,h)),c.has(s[p]))if(c.has(s[h])){const e=d.get(a[u]),r=void 0!==e?o[e]:null;if(null===r){const e=pe(t,o[p]);he(e,n[u]),l[u]=e}else l[u]=he(r,n[u]),pe(t,o[p],r),o[e]=null;u++}else me(o[h]),h--;else me(o[p]),p++;for(;u<=m;){const e=pe(t,l[m+1]);he(e,n[u]),l[u++]=e}for(;p<=h;){const t=o[p++];null!==t&&me(t)}return this.ut=a,((t,e=ue)=>{t._$AH=e})(t,l),Nt}}),ze=ge(class extends ye{constructor(t){if(super(t),t.type!==fe||"class"!==t.name||t.strings?.length>2)throw Error("`classMap()` can only be used in the `class` attribute and must be the only part in the attribute.")}render(t){return" "+Object.keys(t).filter((e=>t[e])).join(" ")+" "}update(t,[e]){if(void 0===this.st){this.st=new Set,void 0!==t.strings&&(this.nt=new Set(t.strings.join(" ").split(/\s/).filter((t=>""!==t))));for(const t in e)e[t]&&!this.nt?.has(t)&&this.st.add(t);return this.render(e)}const r=t.element.classList;for(const t of this.st)t in e||(r.remove(t),this.st.delete(t));for(const t in e){const i=!!e[t];i===this.st.has(t)||this.nt?.has(t)||(i?(r.add(t),this.st.add(t)):(r.remove(t),this.st.delete(t)))}return Nt}});
/**
 * @license
 * Copyright 2018 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */var Me=function(t,e,r,i){var o,n=arguments.length,a=n<3?e:null===i?i=Object.getOwnPropertyDescriptor(e,r):i;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(t,e,r,i);else for(var s=t.length-1;s>=0;s--)(o=t[s])&&(a=(n<3?o(a):n>3?o(e,r,a):o(e,r))||a);return n>3&&a&&Object.defineProperty(e,r,a),a};let Te=class extends ee{constructor(){super(...arguments),this.width=20,this.enabled=!1,this.label="Click to check",this.userSelectable=!0}handleClick(){this.enabled=!this.enabled;const t=new CustomEvent("checked",{detail:{status:this.enabled}});this.dispatchEvent(t)}get value(){return this.enabled?"true":"false"}set value(t){this.enabled="true"===t}render(){return It`<button
      @click=${this.handleClick}
      aria-label="${this.label}"
      style="height: ${this.width}px;"
      ?disabled=${!this.userSelectable}
    >
      <svg
        width="${this.width}"
        height="${this.width}"
        class=${ze({active:this.enabled})}
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
    </button>`}};Te.styles=[rt`
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
    `],Me([ae()],Te.prototype,"width",void 0),Me([ae()],Te.prototype,"enabled",void 0),Me([ae()],Te.prototype,"label",void 0),Me([ae()],Te.prototype,"userSelectable",void 0),Te=Me([ie("radio-input")],Te);var Le=function(t,e,r,i){var o,n=arguments.length,a=n<3?e:null===i?i=Object.getOwnPropertyDescriptor(e,r):i;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(t,e,r,i);else for(var s=t.length-1;s>=0;s--)(o=t[s])&&(a=(n<3?o(a):n>3?o(e,r,a):o(e,r))||a);return n>3&&a&&Object.defineProperty(e,r,a),a};let je=class extends ee{constructor(){super(...arguments),this.default="",this.activeChoice="",this.hoveringChoiceKey="",this.choices=[{key:"small",title:"Small",description:"Small T-Shirt Size"},{key:"medium",title:"Medium",description:"Medium T-Shirt Size"},{key:"large",title:"Large",description:"A big big t shirt"}]}get value(){return this.activeChoice}set value(t){this.activeChoice=t}handleClick(t){this.activeChoice=t.key;const e=new CustomEvent("selected",{detail:{key:t.key}});this.dispatchEvent(e)}render(){return It`<div id="container">
      ${Se(this.choices,(t=>t.key),(t=>It`<div
            class="choice"
            @click=${()=>{this.handleClick(t)}}
            @mouseenter=${()=>{this.hoveringChoiceKey=t.key}}
            @mouseleave=${()=>{this.hoveringChoiceKey=""}}
          >
            <div>
              <radio-input
                class="${ze({"emulate-hover":this.hoveringChoiceKey===t.key})}"
                .label=${`${t.key}: ${t.description}`}
                .enabled=${this.activeChoice===t.key||""===this.activeChoice&&this.default===t.key}
              ></radio-input>
              <p id="title">${t.title}</p>
            </div>
            ${t.description?It`<p id="description">${t.description}</p>`:""}
          </div>`))}
    </div>`}};je.styles=[rt`
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
    `],Le([ae()],je.prototype,"default",void 0),Le([se()],je.prototype,"activeChoice",void 0),Le([se()],je.prototype,"hoveringChoiceKey",void 0),Le([ae()],je.prototype,"choices",void 0),je=Le([ie("radio-selector")],je);var Be=function(t,e,r,i){var o,n=arguments.length,a=n<3?e:null===i?i=Object.getOwnPropertyDescriptor(e,r):i;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(t,e,r,i);else for(var s=t.length-1;s>=0;s--)(o=t[s])&&(a=(n<3?o(a):n>3?o(e,r,a):o(e,r))||a);return n>3&&a&&Object.defineProperty(e,r,a),a};let Ie=class extends ee{constructor(){super(...arguments),this.open=!1,this.wantedWidth="24rem",this.modaltitle="Needs Title",this.lastFocus=null,this._handleKeyDown=t=>{"Escape"===t.key&&this.open&&this.close()}}close(){this.dispatchEvent(new CustomEvent("close")),this.lastFocus&&this.lastFocus.focus()}connectedCallback(){super.connectedCallback(),window.addEventListener("keydown",this._handleKeyDown)}disconnectedCallback(){window.removeEventListener("keydown",this._handleKeyDown),super.disconnectedCallback()}updated(t){var e;if(super.updated(t),t.has("open")&&this.open&&t.get("open")!==this.open){const t=null===(e=this.shadowRoot)||void 0===e?void 0:e.getElementById("content");t&&(this.lastFocus=document.activeElement,t.focus())}t.has("open")&&!this.open&&t.get("open")!==this.open&&this.lastFocus&&this.lastFocus.focus()}render(){return this.open?It`
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
    `:It``}};Ie.styles=rt`
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
  `,Be([ae({type:Boolean})],Ie.prototype,"open",void 0),Be([ae({type:String})],Ie.prototype,"wantedWidth",void 0),Be([ae()],Ie.prototype,"modaltitle",void 0),Ie=Be([ie("modal-component")],Ie);var Ne=function(t,e,r,i){var o,n=arguments.length,a=n<3?e:null===i?i=Object.getOwnPropertyDescriptor(e,r):i;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(t,e,r,i);else for(var s=t.length-1;s>=0;s--)(o=t[s])&&(a=(n<3?o(a):n>3?o(e,r,a):o(e,r))||a);return n>3&&a&&Object.defineProperty(e,r,a),a};let De=class extends ee{constructor(){super(...arguments),this.disabled=!1,this.expectLoad=!1,this.loading=!1,this.loadingText=""}render(){return It`<button
      ?disabled=${this.disabled||this.loading}
      class=${ze({loading:this.loading})}
      @click=${()=>{this.dispatchEvent(new Event("fl-click"))}}
    >
      ${this.expectLoad?It`
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
          `:It``}
      ${this.loading&&""!==this.loadingText?this.loadingText:It`<slot></slot>`}
    </button>`}};De.styles=rt`
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
  `,Ne([ae()],De.prototype,"disabled",void 0),Ne([ae()],De.prototype,"expectLoad",void 0),Ne([ae()],De.prototype,"loading",void 0),Ne([ae()],De.prototype,"loadingText",void 0),De=Ne([ie("button-component")],De);const Ue=rt`
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
`;var He=function(t,e,r,i){var o,n=arguments.length,a=n<3?e:null===i?i=Object.getOwnPropertyDescriptor(e,r):i;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(t,e,r,i);else for(var s=t.length-1;s>=0;s--)(o=t[s])&&(a=(n<3?o(a):n>3?o(e,r,a):o(e,r))||a);return n>3&&a&&Object.defineProperty(e,r,a),a};let Ve=class extends ee{constructor(){super(),this.open=!1,this.promptID="",this.promptType="text",this.promptTitle="",this.promptDescription="",this.promptValue="",this.promptPattern="",this.promptPatternError="",this.promptRadioChoices=[],this.promptRadioDefault=void 0,this.inputRef=Ce(),this.openPrompt=t=>{let{detail:e}=t;this.promptID=e.id,this.promptType=e.type,this.promptTitle=e.title,this.promptDescription=e.description,this.promptPattern=e.pattern,this.promptPatternError=e.patternError,this.promptRadioChoices=e.radioOptions,this.promptRadioDefault=e.radioDefaultKey,this.open=!0}}connectedCallback(){super.connectedCallback(),window.addEventListener("fl-prompt",this.openPrompt)}disconnectedCallback(){window.removeEventListener("fl-prompt",this.openPrompt),super.disconnectedCallback()}handleChange(t){this.promptValue=t.target.value}submitPrompt(){window.dispatchEvent(new CustomEvent(`fl-response-${this.promptID}`,{detail:{id:this.promptID,value:this.promptValue}})),this.close()}close(t){t&&window.dispatchEvent(new CustomEvent(`fl-response-${this.promptID}`,{detail:{canceled:!0}})),this.reset()}reset(){var t;this.open=!1,this.promptType="text",this.promptTitle="",this.promptDescription="",this.promptID="",this.promptValue="",this.promptPattern="",this.promptPatternError="",this.promptRadioChoices=void 0,this.promptRadioDefault=void 0,(null===(t=this.shadowRoot)||void 0===t?void 0:t.getElementById("input")).value=""}updated(t){super.updated(t),t.has("open")&&!this.open&&t.get("open")!==this.open&&this.reset()}getInputMode(){switch(this.promptType){case"number":return"numeric";case"decimal":return"decimal";case"tel":return"tel";default:return"text"}}getPromptArea(){var t,e;switch(this.promptType){case"textarea":return It`<textarea
          id="input"
          @input=${this.handleChange}
          placeholder="type something"
          style="flex-grow: 1;"
        ></textarea>`;case"radio":return It` <radio-selector
          id="input"
          style="margin: 0.6rem 0;"
          .default=${this.promptRadioDefault}
          .choices=${this.promptRadioChoices}
          @selected=${t=>{this.promptValue=t.detail.key}}
        ></radio-selector>`;default:return It`<input
            ${Re(this.inputRef)}
            pattern=${(t=>t??Dt)
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
          ${this.promptPatternError&&!(null===(e=null===(t=this.inputRef.value)||void 0===t?void 0:t.validity)||void 0===e?void 0:e.valid)?It` <p class="error">${this.promptPatternError}</p>`:""}`}}render(){var t,e,r;return It`<modal-component
      .open=${this.open}
      .modaltitle="${this.promptTitle}"
      @close=${this.close}
    >
      <div id="feedback">
        <div>
          ${""!==this.promptDescription?It` <p>${this.promptDescription}</p> `:void 0}
          ${this.getPromptArea()}
        </div>

        <button-component
          class="big"
          @fl-click=${this.submitPrompt}
          .disabled=${0===this.promptValue.length||void 0!==(null===(t=this.inputRef.value)||void 0===t?void 0:t.validity)&&!(null===(r=null===(e=this.inputRef.value)||void 0===e?void 0:e.validity)||void 0===r?void 0:r.valid)}
          >Submit</button-component
        >
      </div>
    </modal-component>`}};Ve.styles=[rt`
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
    `,Ue],He([ae()],Ve.prototype,"open",void 0),He([se()],Ve.prototype,"promptID",void 0),He([se()],Ve.prototype,"promptType",void 0),He([se()],Ve.prototype,"promptTitle",void 0),He([se()],Ve.prototype,"promptDescription",void 0),He([se()],Ve.prototype,"promptValue",void 0),He([se()],Ve.prototype,"promptPattern",void 0),He([se()],Ve.prototype,"promptPatternError",void 0),He([se()],Ve.prototype,"promptRadioChoices",void 0),He([se()],Ve.prototype,"promptRadioDefault",void 0),Ve=He([ie("prompt-component")],Ve);var Ke=function(t,e,r,i){var o,n=arguments.length,a=n<3?e:null===i?i=Object.getOwnPropertyDescriptor(e,r):i;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(t,e,r,i);else for(var s=t.length-1;s>=0;s--)(o=t[s])&&(a=(n<3?o(a):n>3?o(e,r,a):o(e,r))||a);return n>3&&a&&Object.defineProperty(e,r,a),a};let Fe=class extends ee{constructor(){super(...arguments),this.open=!1,this.promptTitle="",this.promptDescription="",this.proceedText="Confirm",this.cancelText="Cancel",this.openAlert=t=>{let{detail:e}=t;this.promptTitle=e.title,this.promptDescription=e.description,this.proceedText=e.proceedButton||this.proceedText,this.cancelText=e.cancelButton||this.cancelText,this.open=!0}}connectedCallback(){super.connectedCallback(),window.addEventListener("fl-confirm",this.openAlert)}disconnectedCallback(){window.removeEventListener("fl-confirm",this.openAlert),super.disconnectedCallback()}close(){this.reset()}reset(){this.open=!1,this.promptTitle="",this.promptDescription="",this.proceedText="Confirm",this.cancelText="Cancel"}updated(t){super.updated(t),t.has("open")&&!this.open&&t.get("open")!==this.open&&this.reset()}render(){return It`<modal-component
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
    </modal-component>`}};Fe.styles=[rt`
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
    `],Ke([ae()],Fe.prototype,"open",void 0),Ke([se()],Fe.prototype,"promptTitle",void 0),Ke([se()],Fe.prototype,"promptDescription",void 0),Ke([se()],Fe.prototype,"proceedText",void 0),Ke([se()],Fe.prototype,"cancelText",void 0),Fe=Ke([ie("confirm-component")],Fe);var We=function(t,e,r,i){var o,n=arguments.length,a=n<3?e:null===i?i=Object.getOwnPropertyDescriptor(e,r):i;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(t,e,r,i);else for(var s=t.length-1;s>=0;s--)(o=t[s])&&(a=(n<3?o(a):n>3?o(e,r,a):o(e,r))||a);return n>3&&a&&Object.defineProperty(e,r,a),a};let Ye=class extends ee{constructor(){super(...arguments),this.open=!1,this.promptTitle="",this.promptDescription=void 0,this.acknowledge=void 0,this.openAlert=t=>{let{detail:e}=t;this.promptTitle=e.title,this.promptDescription=e.description,this.acknowledge=e.acknowledgeText,this.open=!0}}connectedCallback(){super.connectedCallback(),window.addEventListener("fl-alert",this.openAlert)}disconnectedCallback(){window.removeEventListener("fl-alert",this.openAlert),super.disconnectedCallback()}close(){this.reset()}reset(){this.open=!1,this.promptTitle="",this.promptDescription="",this.acknowledge=void 0}updated(t){super.updated(t),t.has("open")&&!this.open&&t.get("open")!==this.open&&this.reset()}render(){return It`<modal-component
      .open=${this.open}
      .modaltitle="${this.promptTitle}"
      @close=${this.close}
    >
      <div id="feedback">
        <div>
          ${void 0!==this.promptDescription?It` <p>${this.promptDescription}</p> `:void 0}
        </div>

        <button-component class="big" @fl-click=${this.close}
          >${this.acknowledge?this.acknowledge:"OK"}</button-component
        >
      </div>
    </modal-component>`}};Ye.styles=[rt`
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
    `],We([ae()],Ye.prototype,"open",void 0),We([se()],Ye.prototype,"promptTitle",void 0),We([se()],Ye.prototype,"promptDescription",void 0),We([se()],Ye.prototype,"acknowledge",void 0),Ye=We([ie("alert-component")],Ye);const Ge={activity:It`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 37 37"><path fill-rule="evenodd" d="M8 2h21a6 6 0 0 1 6 6v21a6 6 0 0 1-6 6H8a6 6 0 0 1-6-6V8a6 6 0 0 1 6-6M0 8a8 8 0 0 1 8-8h21a8 8 0 0 1 8 8v21a8 8 0 0 1-8 8H8a8 8 0 0 1-8-8zm8.5 1a1.5 1.5 0 1 0 0 3h21a1.5 1.5 0 0 0 0-3zM7 18.5A1.5 1.5 0 0 1 8.5 17h18a1.5 1.5 0 0 1 0 3h-18A1.5 1.5 0 0 1 7 18.5M8.5 25a1.5 1.5 0 0 0 0 3h20a1.5 1.5 0 0 0 0-3z" class="primary" clip-rule="evenodd"/></svg>`,alert:It`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 100 100"><g clip-path="url(#a)"><path stroke-width="5" d="m52.165 17.75 34.641 60c.962 1.667-.24 3.75-2.165 3.75H15.359c-1.924 0-3.127-2.083-2.165-3.75l34.64-60c.963-1.667 3.369-1.667 4.331 0" class="primary-stroke"/><path d="M44.414 40.384A5 5 0 0 1 49.4 35h1.202a5 5 0 0 1 4.985 5.383l-1.114 14.475a4.486 4.486 0 0 1-8.945 0z" class="primary"/><circle cx="50" cy="68" r="5" class="primary"/></g><defs><clipPath id="a"><path fill="#fff" d="M0 0h100v100H0z"/></clipPath></defs></svg>`,check:It`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 17 15"><path stroke-width="2" d="m1 9.5 3.695 3.695a1 1 0 0 0 1.5-.098L15.5 1" class="primary-stroke"/></svg>`,checkmark:It`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 33 33"><circle cx="16.15" cy="16.15" r="16.15" class="primary"/><path stroke="#fff" stroke-width="3" d="m8.604 18.867 3.328 3.328a1 1 0 0 0 1.452-.04L24.3 9.962"/></svg>`,clock:It`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 100 100"><path stroke-width="6" d="M68.5 14.526A39.8 39.8 0 0 0 50 10c-22.091 0-40 17.909-40 40s17.909 40 40 40 40-17.909 40-40c0-8.127-1.336-14.688-5.5-21" class="secondary-stroke"/><path d="M87.255 18.607a5 5 0 1 0-7.071-7.071L45.536 46.184a5 5 0 1 0 7.07 7.07zM24.16 82.33a5 5 0 0 0-8.66-5l-5 8.66a5 5 0 1 0 8.66 5zm51.34 0a5 5 0 1 1 8.66-5l5 8.66a5 5 0 0 1-8.66 5z" class="primary"/></svg>`,cog:It`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 95 95"><path fill-rule="evenodd" d="M43 0a5 5 0 0 0-5 5v8.286c0 1.856-1.237 3.473-2.951 4.185-1.715.712-3.71.432-5.024-.88l-5.86-5.86a5 5 0 0 0-7.07 0l-6.365 6.363a5 5 0 0 0 0 7.072l5.86 5.86c1.313 1.312 1.593 3.308.88 5.023C16.76 36.763 15.143 38 13.287 38H5a5 5 0 0 0-5 5v9a5 5 0 0 0 5 5h8.286c1.856 0 3.473 1.237 4.185 2.951.712 1.715.432 3.71-.88 5.024l-5.86 5.86a5 5 0 0 0 0 7.07l6.363 6.364a5 5 0 0 0 7.072 0l5.86-5.86c1.312-1.312 3.308-1.592 5.023-.88S38 79.858 38 81.714V90a5 5 0 0 0 5 5h9a5 5 0 0 0 5-5v-8.286c0-1.856 1.237-3.473 2.951-4.185 1.715-.712 3.71-.432 5.024.88l5.86 5.86a5 5 0 0 0 7.07 0l6.365-6.363a5 5 0 0 0 0-7.071l-5.86-5.86c-1.313-1.313-1.593-3.308-.88-5.024.71-1.714 2.327-2.951 4.183-2.951H90a5 5 0 0 0 5-5v-9a5 5 0 0 0-5-5h-8.286c-1.856 0-3.473-1.237-4.185-2.951-.712-1.715-.432-3.71.88-5.024l5.86-5.86a5 5 0 0 0 0-7.07l-6.363-6.365a5 5 0 0 0-7.071 0l-5.86 5.86c-1.313 1.313-3.308 1.593-5.024.88C58.237 16.76 57 15.143 57 13.287V5a5 5 0 0 0-5-5zm4 62c8.284 0 15-6.716 15-15s-6.716-15-15-15-15 6.716-15 15 6.716 15 15 15" class="primary" clip-rule="evenodd"/></svg>`,email:It`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 100 100"><path d="M48.19 50.952 7 30a6 6 0 0 1 6-6h74a6 6 0 0 1 6 6L52.765 50.93a5 5 0 0 1-4.574.022" class="primary"/><path fill-rule="evenodd" d="M88 26H12a4 4 0 0 0-4 4v41a4 4 0 0 0 4 4h76a4 4 0 0 0 4-4V30a4 4 0 0 0-4-4m-76-4a8 8 0 0 0-8 8v41a8 8 0 0 0 8 8h76a8 8 0 0 0 8-8V30a8 8 0 0 0-8-8z" class="secondary" clip-rule="evenodd"/></svg>`,flag:It`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 9 13"><path d="M8 5V1l-1.175.294a10 10 0 0 1-5.588-.215L1 1v4l.237.08a10 10 0 0 0 5.588.214z" class="secondary"/><path d="M1 12.5V5m0 0V1l.237.08a10 10 0 0 0 5.588.214L8 1v4l-1.175.294a10 10 0 0 1-5.588-.215z" class="primary-stroke"/></svg>`,home:It`<svg xmlns="http://www.w3.org/2000/svg" class="icon-home" viewBox="0 0 24 24"><path d="M9 22H5a1 1 0 0 1-1-1V11l8-8 8 8v10a1 1 0 0 1-1 1h-4a1 1 0 0 1-1-1v-4a1 1 0 0 0-1-1h-2a1 1 0 0 0-1 1v4a1 1 0 0 1-1 1m3-9a2 2 0 1 0 0-4 2 2 0 0 0 0 4" class="primary"/><path d="m12.01 4.42-8.3 8.3a1 1 0 1 1-1.42-1.41l9.02-9.02a1 1 0 0 1 1.41 0l8.99 9.02a1 1 0 0 1-1.42 1.41z" class="secondary"/></svg>`,info:It`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"><circle cx="12" cy="12" r="11.5" stroke="#fff"/><path stroke-width="2" d="M13.5 18.5V13a1 1 0 0 0-1-1H10m3.5 6.5h-4m4 0h3" class="primary-stroke"/><circle cx="12.5" cy="7" r="2" class="primary"/></svg>`,note:It`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 100 100"><g stroke-width="6" clip-path="url(#a)"><path d="M58.657 3H18C9.716 3 3 9.716 3 18v64c0 8.284 6.716 15 15 15h64c8.284 0 15-6.716 15-15V34.629" class="primary-stroke"/><path d="M48.93 54.861 79.801 3.473a1 1 0 0 1 1.358-.35L92.707 9.79a1 1 0 0 1 .406 1.29l-.049.091L62.38 62.25 42.86 76.275z" class="secondary-stroke"/></g><defs><clipPath id="a"><path fill="#fff" d="M0 0h100v100H0z"/></clipPath></defs></svg>`,"person-group":It`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M12 13a3 3 0 0 1 3-3h4a3 3 0 0 1 3 3v3a1 1 0 0 1-1 1h-8a1 1 0 0 1-1-1 1 1 0 0 1-1 1H3a1 1 0 0 1-1-1v-3a3 3 0 0 1 3-3h4a3 3 0 0 1 3 3M7 9a3 3 0 1 1 0-6 3 3 0 0 1 0 6m10 0a3 3 0 1 1 0-6 3 3 0 0 1 0 6" class="secondary"/><path d="M12 13a3 3 0 1 1 0-6 3 3 0 0 1 0 6m-3 1h6a3 3 0 0 1 3 3v3a1 1 0 0 1-1 1H7a1 1 0 0 1-1-1v-3a3 3 0 0 1 3-3" class="primary"/></svg>`,"person-outline":It`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 30 34"><path fill-rule="evenodd" d="M20 7a5 5 0 1 1-10 0 5 5 0 0 1 10 0m2 0A7 7 0 1 1 8 7a7 7 0 0 1 14 0M.187 21.912C-.438 17.746 2.788 14 7 14l1.694 1.376a10 10 0 0 0 12.612 0L23 14c4.212 0 7.438 3.746 6.813 7.912L28 34H2zm22.38-4.983 1.09-.886a4.89 4.89 0 0 1 4.178 5.572L26.278 32H3.722L2.165 21.615a4.89 4.89 0 0 1 4.178-5.572l1.09.886a12 12 0 0 0 15.134 0" class="primary" clip-rule="evenodd"/></svg>`,person:It`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 30 34"><g clip-path="url(#a)"><path fill-rule="evenodd" d="M20 7a5 5 0 1 1-10 0 5 5 0 0 1 10 0m2 0A7 7 0 1 1 8 7a7 7 0 0 1 14 0M.187 21.912C-.438 17.746 2.788 14 7 14l1.694 1.376a10 10 0 0 0 12.612 0L23 14c4.212 0 7.438 3.746 6.813 7.912L28 34H2z" class="primary" clip-rule="evenodd"/></g><defs><clipPath id="a"><path fill="#fff" d="M0 0h30v34H0z"/></clipPath></defs></svg>`,"phone-disabled":It`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 100 100"><g clip-path="url(#a)"><rect width="31" height="11.499" x="37.69" y="4.483" stroke="#A4A4A4" stroke-width="4" rx="3" transform="rotate(60 37.69 4.483)"/><path stroke="#A4A4A4" stroke-width="4" d="M56.212 88a13 13 0 0 1-9.46-4.082l-.233-.255-14.483-16.209a13 13 0 0 1-2.514-4.191l-.109-.31L20.205 35.7a13 13 0 0 1 1.186-10.876l.196-.315 3.355-5.218 9.737 16.23c.21.348.345.735.4 1.136l.018.174.88 11.26a27 27 0 0 0 12.767 20.893l.719.426 6.43 3.689c.383.22.713.52.965.88l.103.158L65.434 88z"/><rect width="31" height="11.499" x="70.69" y="60.732" stroke="#A4A4A4" stroke-width="4" rx="3" transform="rotate(60 70.69 60.732)"/><circle cx="36.869" cy="60.869" r="19.869" class="primary"/><path fill="#fff" fill-rule="evenodd" d="M30.32 49.486a1 1 0 0 0-1.413 0l-3.683 3.682a1 1 0 0 0 0 1.415l5.908 5.907a1 1 0 0 1 0 1.414l-6.103 6.103a1 1 0 0 0 0 1.414l3.55 3.55a1 1 0 0 0 1.414 0l6.103-6.103a1 1 0 0 1 1.414 0l5.907 5.908a1 1 0 0 0 1.415 0l3.682-3.682a1 1 0 0 0 0-1.415l-5.908-5.907a1 1 0 0 1 0-1.414l6.103-6.103a1 1 0 0 0 0-1.415l-3.55-3.55a1 1 0 0 0-1.413 0l-6.104 6.104a1 1 0 0 1-1.414 0z" clip-rule="evenodd"/></g><defs><clipPath id="a"><path fill="#fff" d="M0 0h100v100H0z"/></clipPath></defs></svg>`,phone:It`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 100 100"><g clip-path="url(#a)"><rect width="35" height="15.499" x="37.422" y="-.249" class="secondary" rx="5" transform="rotate(60 37.422 -.25)"/><path d="M24.13 16.854a1 1 0 0 1 1.698.026l10.566 17.61c.399.664.637 1.412.698 2.184l.88 11.26a25 25 0 0 0 12.486 19.74l6.431 3.689a5 5 0 0 1 1.779 1.73l9.402 15.386A1 1 0 0 1 67.217 90H56.212a15 15 0 0 1-11.185-5.005L30.544 68.787a15 15 0 0 1-3.026-5.193L18.311 36.34a15 15 0 0 1 1.593-12.913z" class="primary"/><rect width="35" height="15.499" x="70.422" y="56" class="secondary" rx="5" transform="rotate(60 70.422 56)"/></g><defs><clipPath id="a"><path fill="#fff" d="M0 0h100v100H0z"/></clipPath></defs></svg>`,pin:It`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 100 100"><path fill="#fff" d="M0 0h100v100H0z"/><path fill-rule="evenodd" d="M34.825 12.39A5 5 0 0 1 39.787 8H60.94a5 5 0 0 1 4.971 4.465l5.939 55.142a5 5 0 0 1-4.971 5.535h-5.264q.036-.456.036-.923c0-2.683-.914-5.153-2.447-7.116A5 5 0 0 0 62.37 59.9l-2.89-26.696a5 5 0 0 0-4.971-4.462H46.4a5 5 0 0 0-4.963 4.386l-3.302 26.697A5 5 0 0 0 41.045 65a11.52 11.52 0 0 0-2.493 8.142h-5.551a5 5 0 0 1-4.963-5.61z" class="primary" clip-rule="evenodd"/><circle cx="49.868" cy="72" r="7" class="secondary"/><rect width="8" height="18" x="46" y="75" class="secondary" rx="3"/></svg>`,search:It`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"><path fill-rule="evenodd" d="M9.92 7.93a5.93 5.93 0 1 1 11.858 0 5.93 5.93 0 0 1-11.859 0M15.848 0a7.93 7.93 0 0 0-6.27 12.785L.293 22.07l1.414 1.414 9.286-9.286A7.93 7.93 0 1 0 15.848 0" class="primary" clip-rule="evenodd"/></svg>`,"sign-out":It`<svg xmlns="http://www.w3.org/2000/svg" class="icon-door-exit" viewBox="0 0 24 24"><path d="M11 4h3a1 1 0 0 1 1 1v3a1 1 0 0 1-2 0V6h-2v12h2v-2a1 1 0 0 1 2 0v3a1 1 0 0 1-1 1h-3v1a1 1 0 0 1-1.27.96l-6.98-2A1 1 0 0 1 2 19V5a1 1 0 0 1 .75-.97l6.98-2A1 1 0 0 1 11 3z" class="primary"/><path d="m18.59 11-1.3-1.3c-.94-.94.47-2.35 1.42-1.4l3 3a1 1 0 0 1 0 1.4l-3 3c-.95.95-2.36-.46-1.42-1.4l1.3-1.3H14a1 1 0 0 1 0-2z" class="secondary"/></svg>`,sort:It`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 76 76"><rect width="70" height="11" x="3" y="16" class="primary" rx="5"/><rect width="62" height="11" x="11" y="33" class="primary" rx="5"/><rect width="54" height="11" x="19" y="50" class="primary" rx="5"/></svg>`,trash:It`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 57 58"><path stroke-width="2" d="M6.13 18.658h44.74a4 4 0 0 1 3.918 4.804l-6.023 29.356a4 4 0 0 1-4.232 3.184L28.97 54.778a6 6 0 0 0-.94 0l-15.563 1.224a4 4 0 0 1-4.232-3.184L2.212 23.462a4 4 0 0 1 3.918-4.805" class="primary-stroke"/><rect width="4" height="30.964" class="secondary" rx="2" transform="matrix(.99209 -.12553 .2006 .97967 9.295 22.952)"/><rect width="4" height="30.964" class="secondary" rx="2" transform="matrix(.9921 .12548 -.20051 .9797 44.157 22.45)"/><rect width="4" height="28.805" x="26.872" y="22.138" class="secondary" rx="2"/><path fill-rule="evenodd" d="M37.036 0a3.68 3.68 0 0 1 3.678 3.679 3.68 3.68 0 0 0 3.679 3.678h9.664a2.943 2.943 0 0 1 0 5.886H2.943a2.943 2.943 0 0 1 0-5.886h9.664a3.68 3.68 0 0 0 3.679-3.678A3.68 3.68 0 0 1 19.964 0zM22.564 2.207a2.207 2.207 0 1 0 0 4.415h11.872a2.207 2.207 0 0 0 0-4.415z" class="primary" clip-rule="evenodd"/></svg>`,unlink:It`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 100 100"><path stroke-width="8" d="m41.035 46-24.49 24.452a5 5 0 0 0 0 7.077l5.957 5.945a5 5 0 0 0 7.065 0L46 67.067M58.195 54l25.103-24.393a5 5 0 0 0-.01-7.183l-6.276-6.06a5 5 0 0 0-6.957.009L53 32.933" class="primary-stroke"/><rect width="8" height="18" x="65" y="74.997" class="shadow" rx="4" transform="rotate(-45 65 74.997)"/><rect width="8" height="18" x="73.498" y="63.489" class="shadow" rx="4" transform="rotate(-75 73.498 63.489)"/><rect width="8" height="18" x="49.681" y="79.357" class="shadow" rx="4" transform="rotate(-15 49.68 79.357)"/><rect width="8" height="18" x="34.445" y="21.543" class="shadow" rx="4" transform="rotate(135 34.445 21.543)"/><rect width="8" height="18" x="24.947" y="33.05" class="shadow" rx="4" transform="rotate(105 24.947 33.05)"/><rect width="8" height="18" x="49.765" y="18.182" class="shadow" rx="4" transform="rotate(165 49.765 18.182)"/></svg>`,"view-hidden":It`<svg xmlns="http://www.w3.org/2000/svg" class="icon-view-hidden" viewBox="0 0 24 24"><path d="M15.1 19.34a8 8 0 0 1-8.86-1.68L1.3 12.7a1 1 0 0 1 0-1.42L4.18 8.4l2.8 2.8a5 5 0 0 0 5.73 5.73l2.4 2.4zM8.84 4.6a8 8 0 0 1 8.7 1.74l4.96 4.95a1 1 0 0 1 0 1.42l-2.78 2.78-2.87-2.87a5 5 0 0 0-5.58-5.58L8.85 4.6z" class="primary"/><path d="m3.3 4.7 16 16a1 1 0 0 0 1.4-1.4l-16-16a1 1 0 0 0-1.4 1.4" class="secondary"/></svg>`,"view-visible":It`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M17.56 17.66a8 8 0 0 1-11.32 0L1.3 12.7a1 1 0 0 1 0-1.42l4.95-4.95a8 8 0 0 1 11.32 0l4.95 4.95a1 1 0 0 1 0 1.42l-4.95 4.95zM11.9 17a5 5 0 1 0 0-10 5 5 0 0 0 0 10" class="primary"/><circle cx="12" cy="12" r="3" class="secondary"/></svg>`,xmark:It`<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 33 33"><circle cx="16.149" cy="16.149" r="16.149" class="primary"/><path stroke="#fff" stroke-width="3" d="m9.81 9.96 6.34 6.34m6.339 6.339-6.34-6.339m0 0 6.34-6.34m-6.34 6.34-6.338 6.339"/></svg>`};var qe,Je=function(t,e,r,i){var o,n=arguments.length,a=n<3?e:null===i?i=Object.getOwnPropertyDescriptor(e,r):i;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(t,e,r,i);else for(var s=t.length-1;s>=0;s--)(o=t[s])&&(a=(n<3?o(a):n>3?o(e,r,a):o(e,r))||a);return n>3&&a&&Object.defineProperty(e,r,a),a};let Qe=qe=class extends ee{constructor(){super(...arguments),this.name="info",this.size="24px",this.hoverable=!1,this.colorway="primary"}render(){var t;const e=null!==(t=qe.colorways[this.colorway])&&void 0!==t?t:qe.colorways.primary;return It`
      <div
        class=${ze({hoverable:this.hoverable})}
        style="
          --size: ${this.size};
          --primary: ${e.primary};
          --secondary: ${e.secondary};
          --shadow: ${e.shadow};
        "
      >
        ${Ge[this.name]}
      </div>
    `}};function Ze(t){return"function"==typeof t?t():t}Qe.colorways={primary:{primary:"var(--primary-600)",secondary:"var(--primary-500, #327eff)",shadow:"var(--gray-400, #989898)"},danger:{primary:"var(--danger-600, red)",secondary:"var(--danger-500, pink)",shadow:"var(--gray-500, #888)"}},Qe.styles=rt`
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
  `,Je([ae()],Qe.prototype,"name",void 0),Je([ae()],Qe.prototype,"size",void 0),Je([ae({type:Boolean})],Qe.prototype,"hoverable",void 0),Je([ae()],Qe.prototype,"colorway",void 0),Qe=qe=Je([ie("ui-icon")],Qe);class Xe extends Event{static{this.eventName="lit-state-changed"}constructor(t,e,r){super(Xe.eventName,{cancelable:!1}),this.key=t,this.value=e,this.state=r}}const tr=(t,e)=>e!==t&&(e==e||t==t);class er extends EventTarget{static{this.finalized=!1}static initPropertyMap(){this.propertyMap||(this.propertyMap=new Map)}get propertyMap(){return this.constructor.propertyMap}get stateValue(){return Object.fromEntries([...this.propertyMap].map((([t])=>[t,this[t]])))}constructor(){super(),this.hookMap=new Map,this.constructor.finalize(),this.propertyMap&&[...this.propertyMap].forEach((([t,e])=>{if(void 0!==e.initialValue){const r=Ze(e.initialValue);this[t]=r,e.value=r}}))}static finalize(){if(this.finalized)return!1;this.finalized=!0;const t=Object.keys(this.properties||{});for(const e of t)this.createProperty(e,this.properties[e]);return!0}static createProperty(t,e){this.finalize();const r="symbol"==typeof t?Symbol():`__${t}`,i=this.getPropertyDescriptor(String(t),r,e);Object.defineProperty(this.prototype,t,i)}static getPropertyDescriptor(t,e,r){const i=r?.hasChanged||tr;return{get(){return this[e]},set(r){const o=this[t];this[e]=r,!0===i(r,o)&&this.dispatchStateEvent(t,r,this)},configurable:!0,enumerable:!0}}reset(){this.hookMap.forEach((t=>t.reset())),[...this.propertyMap].filter((([t,e])=>!(!0===e.skipReset||void 0===e.resetValue))).forEach((([t,e])=>{this[t]=e.resetValue}))}subscribe(t,e,r){e&&!Array.isArray(e)&&(e=[e]);const i=r=>{e&&!e.includes(r.key)||t(r.key,r.value,this)};return this.addEventListener(Xe.eventName,i,r),()=>this.removeEventListener(Xe.eventName,i)}dispatchStateEvent(t,e,r){this.dispatchEvent(new Xe(t,e,r))}}class rr{constructor(t,e,r){this.host=t,this.state=e,this.callback=r||(()=>this.host.requestUpdate()),this.host.addController(this)}hostConnected(){this.state.addEventListener(Xe.eventName,this.callback),this.callback()}hostDisconnected(){this.state.removeEventListener(Xe.eventName,this.callback)}}function ir(t){return(e,r)=>{if(Object.getOwnPropertyDescriptor(e,r))throw new Error("@property must be called before all state decorators");const i=e.constructor;i.initPropertyMap();const o=e.hasOwnProperty(r);return i.propertyMap.set(r,{...t,initialValue:t?.value,resetValue:t?.value}),i.createProperty(r,t),o?Object.getOwnPropertyDescriptor(e,r):void 0}}new URL(location.href);const or={prefix:"_ls"};function nr(t){return t={...or,...t},(e,r)=>{const i=Object.getOwnPropertyDescriptor(e,r);if(!i)throw new Error("@local-storage decorator need to be called after @property");const o=`${t?.prefix||""}_${t?.key||String(r)}`,n=e.constructor,a=n.propertyMap.get(r),s=a?.type;if(a){const e=a.initialValue;a.initialValue=()=>function(t,e){if(null!==t&&(e===Boolean||e===Number||e===Array||e===Object))try{t=JSON.parse(t)}catch(e){console.warn("cannot parse value",t)}return t}(localStorage.getItem(o),s)??Ze(e),n.propertyMap.set(r,{...a,...t})}const l=i?.set,c={...i,set:function(t){void 0!==t&&localStorage.setItem(o,s===Object||s===Array?JSON.stringify(t):t),l&&l.call(this,t)}};Object.defineProperty(n.prototype,r,c)}}var ar=function(t,e,r,i){var o,n=arguments.length,a=n<3?e:null===i?i=Object.getOwnPropertyDescriptor(e,r):i;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(t,e,r,i);else for(var s=t.length-1;s>=0;s--)(o=t[s])&&(a=(n<3?o(a):n>3?o(e,r,a):o(e,r))||a);return n>3&&a&&Object.defineProperty(e,r,a),a};class sr extends er{get info(){return this.infoRaw?JSON.parse(this.infoRaw):void 0}set info(t){this.infoRaw=t?JSON.stringify(t):""}get username(){return this.info.username.split("@")[0]}get email(){return this.info.email}get role(){return this.info.role}set extras(t){this.rawExtras=JSON.stringify(t)}get extras(){if(""===this.rawExtras||void 0===this.rawExtras)return{};try{return JSON.parse(this.rawExtras)}catch(t){return console.error(t),{}}}async loadIfNeeded(t=!1){var e;const r=Date.now(),i=""!==this.infoRaw,o=(null===(e=this.permissions)||void 0===e?void 0:e.length)>0,n=0!=+this.loadedAt&&r-+this.loadedAt>9e5;if(!i||!o||n||t)try{const t=await fetch("/api/management/me");if(!t.ok)throw new Error("Failed to fetch /me");const e=await t.json();if(this.info=e.info,this.permissions=e.permissions.join(","),this.loadedAt=`${r}`,sr.GetExtraAboutMe){const t=await sr.GetExtraAboutMe({id:e.info.id,username:e.info.username,role:e.info.role,permissions:e.info.permissions});this.extras=t}}catch(t){console.error("Failed to load /me:",t),window.location.href="/login"}}hasPermission(t){return(this.permissions||"").split(",").includes(t)}hasRole(t){return Array.isArray(t)?t.some((t=>this.role===t)):this.role===t}get isLaunchpad(){const t=document.cookie.split(";");for(let e of t)if(e=e.trim(),e.startsWith("LaunchpadUser="))return e.substring(14)}clear(){localStorage.removeItem("_identity_i"),localStorage.removeItem("_identity_p"),localStorage.removeItem("_identity_la"),localStorage.removeItem("_identity_x")}signOut(){void 0!==sr.SignOutCallback&&sr.SignOutCallback({id:this.info.id,username:this.info.username,role:this.info.role,permissions:this.permissions.split(",")}),this.clear(),window.location.href="/sign-out"}}sr.SignOutCallback=void 0,sr.GetExtraAboutMe=void 0,ar([nr({key:"i",prefix:"_identity"}),ir()],sr.prototype,"infoRaw",void 0),ar([nr({key:"x",prefix:"_identity"}),ir()],sr.prototype,"rawExtras",void 0),ar([nr({key:"p",prefix:"_identity"}),ir()],sr.prototype,"permissions",void 0),ar([nr({key:"la",prefix:"_identity"}),ir()],sr.prototype,"loadedAt",void 0);const lr=new sr;var cr=function(t,e,r,i){var o,n=arguments.length,a=n<3?e:null===i?i=Object.getOwnPropertyDescriptor(e,r):i;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(t,e,r,i);else for(var s=t.length-1;s>=0;s--)(o=t[s])&&(a=(n<3?o(a):n>3?o(e,r,a):o(e,r))||a);return n>3&&a&&Object.defineProperty(e,r,a),a};let dr=class extends ee{constructor(){super(...arguments),this.location={vertical:"top",horizontal:"right"},this.open=!1,this.launchpadForceClosed=!1,this.aboutMeState=new rr(this,lr),this.manageAccountButton=Ce(),this.handleEscape=t=>{"Escape"===t.key&&(this.open=!1)},this.listenForOOBClicks=t=>{if(!this.shadowRoot)return;t.composedPath().includes(this.shadowRoot.host)||(this.open=!1,window.removeEventListener("click",this.listenForOOBClicks))}}updated(){this.setAttribute("location-vertical",this.location.vertical),this.setAttribute("location-horizontal",this.location.horizontal)}connectedCallback(){super.connectedCallback(),lr.loadIfNeeded()}openClicked(){this.open=!this.open,setTimeout((()=>{this.open?(window.addEventListener("click",this.listenForOOBClicks),window.addEventListener("keydown",this.handleEscape),this.updateComplete.then((()=>{setTimeout((()=>{var t;null===(t=this.manageAccountButton.value)||void 0===t||t.focus()}),100)}))):(window.removeEventListener("click",this.listenForOOBClicks),window.removeEventListener("keydown",this.handleEscape))}))}render(){return It` <div id="container">
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
              ${Re(this.manageAccountButton)}
              @click=${()=>{window.location.href="/profile"}}
            >
              <ui-icon name="cog" size="1rem"></ui-icon>
              Manage Account
              <span></span>
            </button>
            ${lr.hasPermission("view.ls-admin")?It`
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

      ${void 0===lr.isLaunchpad||this.launchpadForceClosed?void 0:It`<div id="launchpad">
            <div>
              <p id="launchpad-status">Launchpad</p>
              <p id="launchpad-user">Viewing app as ${lr.isLaunchpad}</p>
            </div>

            <div>
              <button
                @click=${()=>{this.launchpadForceClosed=!0}}
              >
                &times;
              </button>
            </div>
          </div>`}`}};dr.styles=[rt`
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
    `],cr([ae()],dr.prototype,"location",void 0),cr([se()],dr.prototype,"open",void 0),cr([se()],dr.prototype,"launchpadForceClosed",void 0),dr=cr([ie("locksmith-user-icon")],dr);var pr=function(t,e,r,i){var o,n=arguments.length,a=n<3?e:null===i?i=Object.getOwnPropertyDescriptor(e,r):i;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(t,e,r,i);else for(var s=t.length-1;s>=0;s--)(o=t[s])&&(a=(n<3?o(a):n>3?o(e,r,a):o(e,r))||a);return n>3&&a&&Object.defineProperty(e,r,a),a};let hr=0,ur=class extends ee{constructor(){super(),this.toasts=[],window.addEventListener("close-toast",(t=>{const e=t.detail.id;this.removeToastById(e)})),window.addEventListener("do-toast",(t=>{const{id:e,text:r,danger:i,persist:o,actionText:n,onClick:a,duration:s}=t.detail,l=void 0!==s?s:i?8e3:2500,c=null!=e?e:hr++,d={id:c,text:r,danger:i,duration:l,persist:o,actionText:n,onClick:a};this.toasts=[...this.toasts,d],o||setTimeout((()=>this.removeToastById(c)),l)}))}removeToastById(t){const e=this.toasts.find((e=>e.id===t));e&&(e.removing=!0,this.toasts=[...this.toasts],setTimeout((()=>{this.toasts=this.toasts.filter((e=>e.id!==t))}),200))}render(){return It`<div id="root">
      ${Se(this.toasts,(t=>t.id),(t=>{var e;return It`
          <div
            class=${ze({toast:!0,danger:t.danger,"slide-out":null!==(e=t.removing)&&void 0!==e&&e})}
            style="--toast-duration: ${t.duration}ms"
          >
            ${t.text}
            ${t.actionText&&t.onClick?It`<button
                  class="action-btn"
                  @click=${()=>{var e;null===(e=t.onClick)||void 0===e||e.call(t),this.removeToastById(t.id)}}
                >
                  ${t.actionText}
                </button>`:null}
            ${t.persist&&void 0===t.onClick?It`<button
                  class="close-btn"
                  @click=${()=>this.removeToastById(t.id)}
                >
                  &times;
                </button>`:null}

            <div
              class="progress-bar"
              style=${t.persist?"display: none":`animation: shrink ${t.duration}ms linear forwards`}
            ></div>
          </div>
        `}))}
    </div>`}};function mr(t){var e,r;window.dispatchEvent(new CustomEvent("do-toast",{detail:{id:t.id,text:t.text,danger:null!==(e=t.danger)&&void 0!==e&&e,persist:null!==(r=t.persist)&&void 0!==r&&r,actionText:t.actionText,onClick:t.onClick,duration:t.duration}}))}ur.styles=[rt`
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
    `],pr([se()],ur.prototype,"toasts",void 0),ur=pr([ie("toast-component")],ur);var fr=function(t,e,r,i){var o,n=arguments.length,a=n<3?e:null===i?i=Object.getOwnPropertyDescriptor(e,r):i;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(t,e,r,i);else for(var s=t.length-1;s>=0;s--)(o=t[s])&&(a=(n<3?o(a):n>3?o(e,r,a):o(e,r))||a);return n>3&&a&&Object.defineProperty(e,r,a),a};let vr=class extends ee{constructor(){super(),this.locale="en",this.mobileNavOpen=!1,this.includeFooter=!0,this.navbars={};this.navbars={user:{left:[{name:"Overview",path:"/locksmith"},{name:"Users",path:"/locksmith/users",dropdown:[]}],right:[]}},window.ononline=this.onlineNotice,window.onoffline=this.offlineNotice}onlineNotice(){var t;console.log("notyfing toast"),t="network",window.dispatchEvent(new CustomEvent("close-toast",{detail:{id:t}})),mr({id:"net",text:"You are back online."})}offlineNotice(){mr({id:"network",text:"You are offline. Please reconnect to WiFi.",persist:!0,danger:!0})}removeTrailingSlash(t){return t.endsWith("/")?t.slice(0,-1):t}renderDropdownNavbar(t){return It`<section class="dropdown">
      <a
        class="${this.removeTrailingSlash(window.location.pathname)===t.path?"active":""}"
        href="${t.path}"
        >${t.name}</a
      >
      ${void 0!==t.dropdown?It`
            <section>
              ${t.dropdown.map((t=>It`<button @click=${()=>J.go(t.path)}>
                    ${t.name}
                  </button>`))}
            </section>
          `:""}
    </section>`}render(){return It`<div id="root">
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
              ${this.navbars.user.left.map((t=>void 0!==t.dropdown?this.renderDropdownNavbar(t):It`<a
                  class="${this.removeTrailingSlash(window.location.pathname)===this.removeTrailingSlash(t.path)?"active":""}"
                  href="${t.path}"
                  >${t.name}</a
                >`))}
            </div>

            <div>
              <!-- <global-search-component></global-search-component> -->

              ${this.navbars.user.right.map((t=>void 0!==t.dropdown?this.renderDropdownNavbar(t):It`<a
                  class="${this.removeTrailingSlash(window.location.pathname)===this.removeTrailingSlash(t.path)?"active":""}"
                  href="${t.path}"
                  >${t.name}</a
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
            ${this.navbars.user.left.map((t=>It`<a
                  class="${window.location.pathname===t.path?"active":""}"
                  href="${t.path}"
                  >${t.name}</a
                >`))}
            ${this.navbars.user.right.map((t=>It`<a
                  class="${window.location.pathname===t.path?"active":""}"
                  href="${t.path}"
                  >${t.name}</a
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
    </div>`}};vr.styles=[le,rt`
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
    `],fr([ae({type:String})],vr.prototype,"locale",void 0),fr([ae({type:Boolean})],vr.prototype,"mobileNavOpen",void 0),fr([ae({type:Boolean})],vr.prototype,"includeFooter",void 0),fr([ae({type:Object})],vr.prototype,"navbars",void 0),vr=fr([ie("locksmith-layout")],vr);var gr=function(t,e,r,i){var o,n=arguments.length,a=n<3?e:null===i?i=Object.getOwnPropertyDescriptor(e,r):i;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(t,e,r,i);else for(var s=t.length-1;s>=0;s--)(o=t[s])&&(a=(n<3?o(a):n>3?o(e,r,a):o(e,r))||a);return n>3&&a&&Object.defineProperty(e,r,a),a};let yr=class extends ee{constructor(){super(),this.urlBase="",this.urlGroup="",this.urlKey="",this.pages={},this.activePageKey="",this.activePageComponent=void 0,this.activePage=void 0,this.listenForPageChanges=t=>{const{location:e}=t.detail,{pathname:r}=e;this.destruct(),0!==r.substring(this.urlBase.length,r.length).length&&(this.urlKey=e.params.urlKey,this.urlGroup=e.params.urlGroup)},window.addEventListener("vaadin-router-location-changed",this.listenForPageChanges),this.addEventListener("destruct",this.destruct)}destruct(){window.removeEventListener("vaadin-router-location-changed",this.listenForPageChanges),this.removeEventListener("destruct",this.destruct)}firstUpdated(){setTimeout((()=>{const t=this.getRouteParams();if(this.urlGroup=t.urlGroup||"",this.urlKey=t.urlKey||"",""===this.activePageKey){if(""===this.urlGroup||""===this.urlKey)return void this.loadPage(Object.keys(this.pages)[0],0);const t=Object.keys(this.pages).find((t=>t.toLowerCase()===this.urlGroup.toString().toLowerCase()));if(t){const e=this.pages[t].find((t=>t.PageKey===this.urlKey));e?this.loadPage(t,e.SortIndex):this.loadPage(Object.keys(this.pages)[0],0)}else this.loadPage(Object.keys(this.pages)[0],0)}}),0)}getRouteParams(){if(Pr){return Pr.location.params}return{}}async loadPage(t,e,r){this.activePageComponent=void 0;const i=this.pages[t].findIndex((t=>t.SortIndex===e));if(-1===i)throw new Error("page index does not exist!");const o=this.pages[t][i];if(this.urlGroup!==t.toLowerCase()&&this.urlKey!==o.PageKey||o.PageKey!==this.activePageKey){const e=`${this.urlBase}/${t.toLowerCase()}/${o.PageKey}`;r||(""===this.activePageKey?window.history.replaceState({group:t,page:o.PageKey},"",e):window.history.pushState({group:t,page:o.PageKey},"",e))}this.activePage=o,this.activePageKey=o.PageKey;const n=new o.PageComponent;if(void 0!==o.LoadProps){const t=o.LoadProps();n.setProps(t)}await n.OnPageLoad(),this.activePageComponent=n}render(){var t,e;return It` <locksmith-layout>
      <aside class="aside" slot="aside">
        ${Object.keys(this.pages).map((t=>It`<div>
              <h3>${t}</h3>
              <div id="buttons">
                ${this.pages[t].sort(((t,e)=>t.SortIndex-e.SortIndex)).map((e=>It`<button
                        @click=${()=>this.loadPage(t,e.SortIndex)}
                        class="${this.activePageKey===e.PageKey?"active":""}"
                      >
                        ${e.PageName}
                      </button>`))}
              </div>
            </div>`))}
        ${void 0!==(null===(t=this.activePageComponent)||void 0===t?void 0:t.ExtraLeftAside)?It`<hr />
              ${this.activePageComponent.ExtraLeftAside()}`:It``}
      </aside>

      <div id="main">
        ${void 0!==this.activePageComponent?It`${this.activePageComponent}`:"Loading.."}
      </div>

      ${void 0!==(null===(e=this.activePageComponent)||void 0===e?void 0:e.GetRightAside)?It`<div class="aside right" slot="aside-right">
            ${this.activePageComponent.GetRightAside()}
          </div>`:It``}
    </locksmith-layout>`}};yr.styles=[le,rt`
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
    `],gr([ae({type:String})],yr.prototype,"urlBase",void 0),gr([ae({type:String})],yr.prototype,"urlGroup",void 0),gr([ae({type:String})],yr.prototype,"urlKey",void 0),gr([ae({type:Array})],yr.prototype,"pages",void 0),gr([se()],yr.prototype,"activePageKey",void 0),gr([se()],yr.prototype,"activePageComponent",void 0),gr([se()],yr.prototype,"activePage",void 0),yr=gr([ie("locksmith-subnav-layout")],yr);var br=function(t,e,r,i){var o,n=arguments.length,a=n<3?e:null===i?i=Object.getOwnPropertyDescriptor(e,r):i;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(t,e,r,i);else for(var s=t.length-1;s>=0;s--)(o=t[s])&&(a=(n<3?o(a):n>3?o(e,r,a):o(e,r))||a);return n>3&&a&&Object.defineProperty(e,r,a),a};let wr=class extends ee{constructor(){super(...arguments),this.users=[],this.loading=!0}OnPageLoad(){this.fetchUsers()}async fetchUsers(){const t=await fetch("/api/users/list"),e=await t.json();this.users=e,this.loading=!1}render(){return It`
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
          ${this.users.map((t=>{var e;return It`
              <button href="#" class="item">
                <p id="title">${null!==(e=t.username)&&void 0!==e?e:t.email}</p>
                <p id="description">
                  ${t.role} &bull; ${t.sessions} active sessions
                </p>
              </button>
            `}))}
        </div>
      </div>
    `}};async function xr(t){const e=crypto.randomUUID();return t.id=e,new Promise((e=>{window.addEventListener(`fl-response-${t.id}`,(t=>{e(t.detail)}),{once:!0}),window.dispatchEvent(new CustomEvent("fl-prompt",{detail:t}))}))}wr.styles=[le,rt``],br([se()],wr.prototype,"users",void 0),br([se()],wr.prototype,"loading",void 0),wr=br([ie("locksmith-users-subpage")],wr);var $r=function(t,e,r,i){var o,n=arguments.length,a=n<3?e:null===i?i=Object.getOwnPropertyDescriptor(e,r):i;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(t,e,r,i);else for(var s=t.length-1;s>=0;s--)(o=t[s])&&(a=(n<3?o(a):n>3?o(e,r,a):o(e,r))||a);return n>3&&a&&Object.defineProperty(e,r,a),a};let kr=class extends ee{constructor(){super(...arguments),this.invites=[],this.loading=!0,this.sendInviteBtnRef=Ce()}OnPageLoad(){this.fetchInvites()}async fetchInvites(){const t=await fetch("/api/users/invitations"),e=await t.json();this.invites=e.reverse(),this.loading=!1}async sendInvite(){const t=await xr({id:"send-invites",title:"Send a new Invite",description:"What is the new users email address? You will be able to select their role on the next screen.",type:"text"});if(t.canceled)return;const e=await xr({id:"invite-role",title:"What is this users role?",description:"This will give them specific permissions, be careful!",type:"radio",radioOptions:[{key:"user",title:"User"},{key:"admin",title:"Admin"}]});if(!e.canceled){this.sendInviteBtnRef.value.loading=!0;try{const i=await async function(t,e){return fetch(t,e)}("/api/users/invite",{method:"POST",body:JSON.stringify({email:t.value,role:e.value})});if(200!==i.status){if(409===i.status)return r={title:"This user has already been invited.",description:"Please use a different email or re-issue the old invite."},void window.dispatchEvent(new CustomEvent("fl-alert",{detail:r}));throw new Error("Got a bad status code: "+i.status)}mr({text:"Invitation email has been sent."}),this.invites=[{email:t.value,role:e.value,inviter:"",sentAt:+new Date/1e3,userid:""},...this.invites]}catch(t){mr({text:`Failed to send invite: ${t.message}`,danger:!0})}finally{this.sendInviteBtnRef.value.loading=!1}var r}}render(){return It`
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
            ${Re(this.sendInviteBtnRef)}
            @fl-click=${()=>{this.sendInvite()}}
          >
            Send Invite
          </button-component>
        </div>
      </header>

      <div class="widget" style="${this.loading?"display: none;":""}">
        <div class="list">
          ${this.invites.map((t=>It`
              <button href="#" class="item">
                <p id="title">${t.email}</p>
                <p id="description">Invited ${function(t){const e=new Date(1e3*t),r=new Date,i=Math.floor((r.getTime()-e.getTime())/1e3);if(i<60)return`${i}s ago`;const o=Math.floor(i/60);if(o<60)return`${o}m ago`;const n=Math.floor(o/60);if(n<24)return`${n}h ago`;const a=Math.floor(n/24);return a<7?`${a}d ago`:e.toLocaleDateString()}(t.sentAt)}</p>
              </button>
            `))}
        </div>
      </div>
    `}};kr.styles=[le,rt``],$r([se()],kr.prototype,"invites",void 0),$r([se()],kr.prototype,"loading",void 0),kr=$r([ie("locksmith-invitations-subpage")],kr);var _r=function(t,e,r,i){var o,n=arguments.length,a=n<3?e:null===i?i=Object.getOwnPropertyDescriptor(e,r):i;if("object"==typeof Reflect&&"function"==typeof Reflect.decorate)a=Reflect.decorate(t,e,r,i);else for(var s=t.length-1;s>=0;s--)(o=t[s])&&(a=(n<3?o(a):n>3?o(e,r,a):o(e,r))||a);return n>3&&a&&Object.defineProperty(e,r,a),a};let Ar=class extends ee{constructor(){super(...arguments),this.pages={Users:[{PageKey:"all",PageComponent:wr,PageName:"Registered Users",SortIndex:0},{PageKey:"invites",PageComponent:kr,PageName:"Invitations",SortIndex:1}]}}onBeforeEnter(t){}render(){return It`
      <locksmith-subnav-layout
        .urlBase=${"/locksmith/users"}
        .pages=${this.pages}
      >
      </locksmith-subnav-layout>
    `}};Ar.styles=[le,rt``],Ar=_r([ie("locksmith-users-page")],Ar),Qe.colorways.primary={primary:"var(--primary-800)",secondary:"var(--primary-300)",shadow:"var(--gray-600)"};const Cr=[{path:"/locksmith/users/:urlGroup?/:urlKey?",component:"locksmith-users-page"},{path:"(.*)",action:t=>{console.log(t),console.warn("Page not found"),J.go("/locksmith/users")},component:"not-found-page"}],Er=document.getElementById("outlet"),Pr=new J(Er);Pr.setRoutes(Cr);export{Pr as router};
//# sourceMappingURL=locksmith-admin.bundle.js.map
