/*! 
 Build based on gin-vue-admin 
 Time : 1682164992000 */
System.register(["./index-legacy.f55fd1af.js"],(function(t,e){"use strict";var r;return{setters:[function(t){t.d1,r=t.d7}],execute:function(){var e,n=t("b",{}),o={},i={},a=[0,26,44,70,100,134,172,196,242,292,346,404,466,532,581,655,733,815,901,991,1085,1156,1258,1364,1474,1588,1706,1828,1921,2051,2185,2323,2465,2611,2761,2876,3034,3196,3362,3532,3706];i.getSymbolSize=function(t){if(!t)throw new Error('"version" cannot be null or undefined');if(t<1||t>40)throw new Error('"version" should be in range from 1 to 40');return 4*t+17},i.getSymbolTotalCodewords=function(t){return a[t]},i.getBCHDigit=function(t){for(var e=0;0!==t;)e++,t>>>=1;return e},i.setToSJISFunction=function(t){if("function"!=typeof t)throw new Error('"toSJISFunc" is not a valid function.');e=t},i.isKanjiModeEnabled=function(){return void 0!==e},i.toSJIS=function(t){return e(t)};var u={};function s(){this.buffer=[],this.length=0}!function(t){t.L={bit:1},t.M={bit:0},t.Q={bit:3},t.H={bit:2},t.isValid=function(t){return t&&void 0!==t.bit&&t.bit>=0&&t.bit<4},t.from=function(e,r){if(t.isValid(e))return e;try{return function(e){if("string"!=typeof e)throw new Error("Param is not a string");switch(e.toLowerCase()){case"l":case"low":return t.L;case"m":case"medium":return t.M;case"q":case"quartile":return t.Q;case"h":case"high":return t.H;default:throw new Error("Unknown EC Level: "+e)}}(e)}catch(n){return r}}}(u),s.prototype={get:function(t){var e=Math.floor(t/8);return 1==(this.buffer[e]>>>7-t%8&1)},put:function(t,e){for(var r=0;r<e;r++)this.putBit(1==(t>>>e-r-1&1))},getLengthInBits:function(){return this.length},putBit:function(t){var e=Math.floor(this.length/8);this.buffer.length<=e&&this.buffer.push(0),t&&(this.buffer[e]|=128>>>this.length%8),this.length++}};var c=s;function f(t){if(!t||t<1)throw new Error("BitMatrix size must be defined and greater than 0");this.size=t,this.data=new Uint8Array(t*t),this.reservedBit=new Uint8Array(t*t)}f.prototype.set=function(t,e,r,n){var o=t*this.size+e;this.data[o]=r,n&&(this.reservedBit[o]=!0)},f.prototype.get=function(t,e){return this.data[t*this.size+e]},f.prototype.xor=function(t,e,r){this.data[t*this.size+e]^=r},f.prototype.isReserved=function(t,e){return this.reservedBit[t*this.size+e]};var l=f,d={};!function(t){var e=i.getSymbolSize;t.getRowColCoords=function(t){if(1===t)return[];for(var r=Math.floor(t/7)+2,n=e(t),o=145===n?26:2*Math.ceil((n-13)/(2*r-2)),i=[n-7],a=1;a<r-1;a++)i[a]=i[a-1]-o;return i.push(6),i.reverse()},t.getPositions=function(e){for(var r=[],n=t.getRowColCoords(e),o=n.length,i=0;i<o;i++)for(var a=0;a<o;a++)0===i&&0===a||0===i&&a===o-1||i===o-1&&0===a||r.push([n[i],n[a]]);return r}}(d);var h={},g=i.getSymbolSize;h.getPositions=function(t){var e=g(t);return[[0,0],[e-7,0],[0,e-7]]};var p={};!function(t){t.Patterns={PATTERN000:0,PATTERN001:1,PATTERN010:2,PATTERN011:3,PATTERN100:4,PATTERN101:5,PATTERN110:6,PATTERN111:7};var e=3,r=3,n=40,o=10;function i(e,r,n){switch(e){case t.Patterns.PATTERN000:return(r+n)%2==0;case t.Patterns.PATTERN001:return r%2==0;case t.Patterns.PATTERN010:return n%3==0;case t.Patterns.PATTERN011:return(r+n)%3==0;case t.Patterns.PATTERN100:return(Math.floor(r/2)+Math.floor(n/3))%2==0;case t.Patterns.PATTERN101:return r*n%2+r*n%3==0;case t.Patterns.PATTERN110:return(r*n%2+r*n%3)%2==0;case t.Patterns.PATTERN111:return(r*n%3+(r+n)%2)%2==0;default:throw new Error("bad maskPattern:"+e)}}t.isValid=function(t){return null!=t&&""!==t&&!isNaN(t)&&t>=0&&t<=7},t.from=function(e){return t.isValid(e)?parseInt(e,10):void 0},t.getPenaltyN1=function(t){for(var r=t.size,n=0,o=0,i=0,a=null,u=null,s=0;s<r;s++){o=i=0,a=u=null;for(var c=0;c<r;c++){var f=t.get(s,c);f===a?o++:(o>=5&&(n+=e+(o-5)),a=f,o=1),(f=t.get(c,s))===u?i++:(i>=5&&(n+=e+(i-5)),u=f,i=1)}o>=5&&(n+=e+(o-5)),i>=5&&(n+=e+(i-5))}return n},t.getPenaltyN2=function(t){for(var e=t.size,n=0,o=0;o<e-1;o++)for(var i=0;i<e-1;i++){var a=t.get(o,i)+t.get(o,i+1)+t.get(o+1,i)+t.get(o+1,i+1);4!==a&&0!==a||n++}return n*r},t.getPenaltyN3=function(t){for(var e=t.size,r=0,o=0,i=0,a=0;a<e;a++){o=i=0;for(var u=0;u<e;u++)o=o<<1&2047|t.get(a,u),u>=10&&(1488===o||93===o)&&r++,i=i<<1&2047|t.get(u,a),u>=10&&(1488===i||93===i)&&r++}return r*n},t.getPenaltyN4=function(t){for(var e=0,r=t.data.length,n=0;n<r;n++)e+=t.data[n];return Math.abs(Math.ceil(100*e/r/5)-10)*o},t.applyMask=function(t,e){for(var r=e.size,n=0;n<r;n++)for(var o=0;o<r;o++)e.isReserved(o,n)||e.xor(o,n,i(t,o,n))},t.getBestMask=function(e,r){for(var n=Object.keys(t.Patterns).length,o=0,i=1/0,a=0;a<n;a++){r(a),t.applyMask(a,e);var u=t.getPenaltyN1(e)+t.getPenaltyN2(e)+t.getPenaltyN3(e)+t.getPenaltyN4(e);t.applyMask(a,e),u<i&&(i=u,o=a)}return o}}(p);var v={},y=u,m=[1,1,1,1,1,1,1,1,1,1,2,2,1,2,2,4,1,2,4,4,2,4,4,4,2,4,6,5,2,4,6,6,2,5,8,8,4,5,8,8,4,5,8,11,4,8,10,11,4,9,12,16,4,9,16,16,6,10,12,18,6,10,17,16,6,11,16,19,6,13,18,21,7,14,21,25,8,16,20,25,8,17,23,25,9,17,23,34,9,18,25,30,10,20,27,32,12,21,29,35,12,23,34,37,12,25,34,40,13,26,35,42,14,28,38,45,15,29,40,48,16,31,43,51,17,33,45,54,18,35,48,57,19,37,51,60,19,38,53,63,20,40,56,66,21,43,59,70,22,45,62,74,24,47,65,77,25,49,68,81],w=[7,10,13,17,10,16,22,28,15,26,36,44,20,36,52,64,26,48,72,88,36,64,96,112,40,72,108,130,48,88,132,156,60,110,160,192,72,130,192,224,80,150,224,264,96,176,260,308,104,198,288,352,120,216,320,384,132,240,360,432,144,280,408,480,168,308,448,532,180,338,504,588,196,364,546,650,224,416,600,700,224,442,644,750,252,476,690,816,270,504,750,900,300,560,810,960,312,588,870,1050,336,644,952,1110,360,700,1020,1200,390,728,1050,1260,420,784,1140,1350,450,812,1200,1440,480,868,1290,1530,510,924,1350,1620,540,980,1440,1710,570,1036,1530,1800,570,1064,1590,1890,600,1120,1680,1980,630,1204,1770,2100,660,1260,1860,2220,720,1316,1950,2310,750,1372,2040,2430];v.getBlocksCount=function(t,e){switch(e){case y.L:return m[4*(t-1)+0];case y.M:return m[4*(t-1)+1];case y.Q:return m[4*(t-1)+2];case y.H:return m[4*(t-1)+3];default:return}},v.getTotalCodewordsCount=function(t,e){switch(e){case y.L:return w[4*(t-1)+0];case y.M:return w[4*(t-1)+1];case y.Q:return w[4*(t-1)+2];case y.H:return w[4*(t-1)+3];default:return}};var E={},b={},A=new Uint8Array(512),C=new Uint8Array(256);!function(){for(var t=1,e=0;e<255;e++)A[e]=t,C[t]=e,256&(t<<=1)&&(t^=285);for(var r=255;r<512;r++)A[r]=A[r-255]}(),b.log=function(t){if(t<1)throw new Error("log("+t+")");return C[t]},b.exp=function(t){return A[t]},b.mul=function(t,e){return 0===t||0===e?0:A[C[t]+C[e]]},function(t){var e=b;t.mul=function(t,r){for(var n=new Uint8Array(t.length+r.length-1),o=0;o<t.length;o++)for(var i=0;i<r.length;i++)n[o+i]^=e.mul(t[o],r[i]);return n},t.mod=function(t,r){for(var n=new Uint8Array(t);n.length-r.length>=0;){for(var o=n[0],i=0;i<r.length;i++)n[i]^=e.mul(r[i],o);for(var a=0;a<n.length&&0===n[a];)a++;n=n.slice(a)}return n},t.generateECPolynomial=function(r){for(var n=new Uint8Array([1]),o=0;o<r;o++)n=t.mul(n,new Uint8Array([1,e.exp(o)]));return n}}(E);var T=E;function M(t){this.genPoly=void 0,this.degree=t,this.degree&&this.initialize(this.degree)}M.prototype.initialize=function(t){this.degree=t,this.genPoly=T.generateECPolynomial(this.degree)},M.prototype.encode=function(t){if(!this.genPoly)throw new Error("Encoder not initialized");var e=new Uint8Array(t.length+this.degree);e.set(t);var r=T.mod(e,this.genPoly),n=this.degree-r.length;if(n>0){var o=new Uint8Array(this.degree);return o.set(r,n),o}return r};var B=M,I={},P={},N={isValid:function(t){return!isNaN(t)&&t>=1&&t<=40}},S={},x="[0-9]+",R="(?:[u3000-u303F]|[u3040-u309F]|[u30A0-u30FF]|[uFF00-uFFEF]|[u4E00-u9FAF]|[u2605-u2606]|[u2190-u2195]|u203B|[u2010u2015u2018u2019u2025u2026u201Cu201Du2225u2260]|[u0391-u0451]|[u00A7u00A8u00B1u00B4u00D7u00F7])+",L="(?:(?![A-Z0-9 $%*+\\-./:]|"+(R=R.replace(/u/g,"\\u"))+")(?:.|[\r\n]))+";S.KANJI=new RegExp(R,"g"),S.BYTE_KANJI=new RegExp("[^A-Z0-9 $%*+\\-./:]+","g"),S.BYTE=new RegExp(L,"g"),S.NUMERIC=new RegExp(x,"g"),S.ALPHANUMERIC=new RegExp("[A-Z $%*+\\-./:]+","g");var U=new RegExp("^"+R+"$"),k=new RegExp("^"+x+"$"),_=new RegExp("^[A-Z0-9 $%*+\\-./:]+$");S.testKanji=function(t){return U.test(t)},S.testNumeric=function(t){return k.test(t)},S.testAlphanumeric=function(t){return _.test(t)},function(t){var e=N,r=S;t.NUMERIC={id:"Numeric",bit:1,ccBits:[10,12,14]},t.ALPHANUMERIC={id:"Alphanumeric",bit:2,ccBits:[9,11,13]},t.BYTE={id:"Byte",bit:4,ccBits:[8,16,16]},t.KANJI={id:"Kanji",bit:8,ccBits:[8,10,12]},t.MIXED={bit:-1},t.getCharCountIndicator=function(t,r){if(!t.ccBits)throw new Error("Invalid mode: "+t);if(!e.isValid(r))throw new Error("Invalid version: "+r);return r>=1&&r<10?t.ccBits[0]:r<27?t.ccBits[1]:t.ccBits[2]},t.getBestModeForData=function(e){return r.testNumeric(e)?t.NUMERIC:r.testAlphanumeric(e)?t.ALPHANUMERIC:r.testKanji(e)?t.KANJI:t.BYTE},t.toString=function(t){if(t&&t.id)return t.id;throw new Error("Invalid mode")},t.isValid=function(t){return t&&t.bit&&t.ccBits},t.from=function(e,r){if(t.isValid(e))return e;try{return function(e){if("string"!=typeof e)throw new Error("Param is not a string");switch(e.toLowerCase()){case"numeric":return t.NUMERIC;case"alphanumeric":return t.ALPHANUMERIC;case"kanji":return t.KANJI;case"byte":return t.BYTE;default:throw new Error("Unknown mode: "+e)}}(e)}catch(n){return r}}}(P),function(t){var e=i,r=v,n=u,o=P,a=N,s=e.getBCHDigit(7973);function c(t,e){return o.getCharCountIndicator(t,e)+4}function f(t,e){var r=0;return t.forEach((function(t){var n=c(t.mode,e);r+=n+t.getBitsLength()})),r}t.from=function(t,e){return a.isValid(t)?parseInt(t,10):e},t.getCapacity=function(t,n,i){if(!a.isValid(t))throw new Error("Invalid QR Code version");void 0===i&&(i=o.BYTE);var u=8*(e.getSymbolTotalCodewords(t)-r.getTotalCodewordsCount(t,n));if(i===o.MIXED)return u;var s=u-c(i,t);switch(i){case o.NUMERIC:return Math.floor(s/10*3);case o.ALPHANUMERIC:return Math.floor(s/11*2);case o.KANJI:return Math.floor(s/13);case o.BYTE:default:return Math.floor(s/8)}},t.getBestVersionForData=function(e,r){var i,a=n.from(r,n.M);if(Array.isArray(e)){if(e.length>1)return function(e,r){for(var n=1;n<=40;n++)if(f(e,n)<=t.getCapacity(n,r,o.MIXED))return n}(e,a);if(0===e.length)return 1;i=e[0]}else i=e;return function(e,r,n){for(var o=1;o<=40;o++)if(r<=t.getCapacity(o,n,e))return o}(i.mode,i.getLength(),a)},t.getEncodedBits=function(t){if(!a.isValid(t)||t<7)throw new Error("Invalid QR Code version");for(var r=t<<12;e.getBCHDigit(r)-s>=0;)r^=7973<<e.getBCHDigit(r)-s;return t<<12|r}}(I);var O={},j=i,F=j.getBCHDigit(1335);O.getEncodedBits=function(t,e){for(var r=t.bit<<3|e,n=r<<10;j.getBCHDigit(n)-F>=0;)n^=1335<<j.getBCHDigit(n)-F;return 21522^(r<<10|n)};var H={},z=P;function D(t){this.mode=z.NUMERIC,this.data=t.toString()}D.getBitsLength=function(t){return 10*Math.floor(t/3)+(t%3?t%3*3+1:0)},D.prototype.getLength=function(){return this.data.length},D.prototype.getBitsLength=function(){return D.getBitsLength(this.data.length)},D.prototype.write=function(t){var e,r,n;for(e=0;e+3<=this.data.length;e+=3)r=this.data.substr(e,3),n=parseInt(r,10),t.put(n,10);var o=this.data.length-e;o>0&&(r=this.data.substr(e),n=parseInt(r,10),t.put(n,3*o+1))};var Y=D,J=P,K=["0","1","2","3","4","5","6","7","8","9","A","B","C","D","E","F","G","H","I","J","K","L","M","N","O","P","Q","R","S","T","U","V","W","X","Y","Z"," ","$","%","*","+","-",".","/",":"];function q(t){this.mode=J.ALPHANUMERIC,this.data=t}q.getBitsLength=function(t){return 11*Math.floor(t/2)+t%2*6},q.prototype.getLength=function(){return this.data.length},q.prototype.getBitsLength=function(){return q.getBitsLength(this.data.length)},q.prototype.write=function(t){var e;for(e=0;e+2<=this.data.length;e+=2){var r=45*K.indexOf(this.data[e]);r+=K.indexOf(this.data[e+1]),t.put(r,11)}this.data.length%2&&t.put(K.indexOf(this.data[e]),6)};var V=q,Q=function(t){for(var e=[],r=t.length,n=0;n<r;n++){var o=t.charCodeAt(n);if(o>=55296&&o<=56319&&r>n+1){var i=t.charCodeAt(n+1);i>=56320&&i<=57343&&(o=1024*(o-55296)+i-56320+65536,n+=1)}o<128?e.push(o):o<2048?(e.push(o>>6|192),e.push(63&o|128)):o<55296||o>=57344&&o<65536?(e.push(o>>12|224),e.push(o>>6&63|128),e.push(63&o|128)):o>=65536&&o<=1114111?(e.push(o>>18|240),e.push(o>>12&63|128),e.push(o>>6&63|128),e.push(63&o|128)):e.push(239,191,189)}return new Uint8Array(e).buffer},$=P;function X(t){this.mode=$.BYTE,"string"==typeof t&&(t=Q(t)),this.data=new Uint8Array(t)}X.getBitsLength=function(t){return 8*t},X.prototype.getLength=function(){return this.data.length},X.prototype.getBitsLength=function(){return X.getBitsLength(this.data.length)},X.prototype.write=function(t){for(var e=0,r=this.data.length;e<r;e++)t.put(this.data[e],8)};var Z=X,W=P,G=i;function tt(t){this.mode=W.KANJI,this.data=t}tt.getBitsLength=function(t){return 13*t},tt.prototype.getLength=function(){return this.data.length},tt.prototype.getBitsLength=function(){return tt.getBitsLength(this.data.length)},tt.prototype.write=function(t){var e;for(e=0;e<this.data.length;e++){var r=G.toSJIS(this.data[e]);if(r>=33088&&r<=40956)r-=33088;else{if(!(r>=57408&&r<=60351))throw new Error("Invalid SJIS character: "+this.data[e]+"\nMake sure your charset is UTF-8");r-=49472}r=192*(r>>>8&255)+(255&r),t.put(r,13)}};var et=tt,rt={exports:{}};!function(t){var e={single_source_shortest_paths:function(t,r,n){var o={},i={};i[r]=0;var a,u,s,c,f,l,d,h=e.PriorityQueue.make();for(h.push(r,0);!h.empty();)for(s in u=(a=h.pop()).value,c=a.cost,f=t[u]||{})f.hasOwnProperty(s)&&(l=c+f[s],d=i[s],(void 0===i[s]||d>l)&&(i[s]=l,h.push(s,l),o[s]=u));if(void 0!==n&&void 0===i[n]){var g=["Could not find a path from ",r," to ",n,"."].join("");throw new Error(g)}return o},extract_shortest_path_from_predecessor_list:function(t,e){for(var r=[],n=e;n;)r.push(n),t[n],n=t[n];return r.reverse(),r},find_path:function(t,r,n){var o=e.single_source_shortest_paths(t,r,n);return e.extract_shortest_path_from_predecessor_list(o,n)},PriorityQueue:{make:function(t){var r,n=e.PriorityQueue,o={};for(r in t=t||{},n)n.hasOwnProperty(r)&&(o[r]=n[r]);return o.queue=[],o.sorter=t.sorter||n.default_sorter,o},default_sorter:function(t,e){return t.cost-e.cost},push:function(t,e){var r={value:t,cost:e};this.queue.push(r),this.queue.sort(this.sorter)},pop:function(){return this.queue.shift()},empty:function(){return 0===this.queue.length}}};t.exports=e}(rt),function(t){var e=P,r=Y,n=V,o=Z,a=et,u=S,s=i,c=rt.exports;function f(t){return unescape(encodeURIComponent(t)).length}function l(t,e,r){for(var n,o=[];null!==(n=t.exec(r));)o.push({data:n[0],index:n.index,mode:e,length:n[0].length});return o}function d(t){var r,n,o=l(u.NUMERIC,e.NUMERIC,t),i=l(u.ALPHANUMERIC,e.ALPHANUMERIC,t);return s.isKanjiModeEnabled()?(r=l(u.BYTE,e.BYTE,t),n=l(u.KANJI,e.KANJI,t)):(r=l(u.BYTE_KANJI,e.BYTE,t),n=[]),o.concat(i,r,n).sort((function(t,e){return t.index-e.index})).map((function(t){return{data:t.data,mode:t.mode,length:t.length}}))}function h(t,i){switch(i){case e.NUMERIC:return r.getBitsLength(t);case e.ALPHANUMERIC:return n.getBitsLength(t);case e.KANJI:return a.getBitsLength(t);case e.BYTE:return o.getBitsLength(t)}}function g(t,i){var u,c=e.getBestModeForData(t);if((u=e.from(i,c))!==e.BYTE&&u.bit<c.bit)throw new Error('"'+t+'" cannot be encoded with mode '+e.toString(u)+".\n Suggested mode is: "+e.toString(c));switch(u!==e.KANJI||s.isKanjiModeEnabled()||(u=e.BYTE),u){case e.NUMERIC:return new r(t);case e.ALPHANUMERIC:return new n(t);case e.KANJI:return new a(t);case e.BYTE:return new o(t)}}t.fromArray=function(t){return t.reduce((function(t,e){return"string"==typeof e?t.push(g(e,null)):e.data&&t.push(g(e.data,e.mode)),t}),[])},t.fromString=function(r,n){for(var o=function(t){for(var r=[],n=0;n<t.length;n++){var o=t[n];switch(o.mode){case e.NUMERIC:r.push([o,{data:o.data,mode:e.ALPHANUMERIC,length:o.length},{data:o.data,mode:e.BYTE,length:o.length}]);break;case e.ALPHANUMERIC:r.push([o,{data:o.data,mode:e.BYTE,length:o.length}]);break;case e.KANJI:r.push([o,{data:o.data,mode:e.BYTE,length:f(o.data)}]);break;case e.BYTE:r.push([{data:o.data,mode:e.BYTE,length:f(o.data)}])}}return r}(d(r,s.isKanjiModeEnabled())),i=function(t,r){for(var n={},o={start:{}},i=["start"],a=0;a<t.length;a++){for(var u=t[a],s=[],c=0;c<u.length;c++){var f=u[c],l=""+a+c;s.push(l),n[l]={node:f,lastCount:0},o[l]={};for(var d=0;d<i.length;d++){var g=i[d];n[g]&&n[g].node.mode===f.mode?(o[g][l]=h(n[g].lastCount+f.length,f.mode)-h(n[g].lastCount,f.mode),n[g].lastCount+=f.length):(n[g]&&(n[g].lastCount=f.length),o[g][l]=h(f.length,f.mode)+4+e.getCharCountIndicator(f.mode,r))}}i=s}for(var p=0;p<i.length;p++)o[i[p]].end=0;return{map:o,table:n}}(o,n),a=c.find_path(i.map,"start","end"),u=[],l=1;l<a.length-1;l++)u.push(i.table[a[l]].node);return t.fromArray(function(t){return t.reduce((function(t,e){var r=t.length-1>=0?t[t.length-1]:null;return r&&r.mode===e.mode?(t[t.length-1].data+=e.data,t):(t.push(e),t)}),[])}(u))},t.rawSplit=function(e){return t.fromArray(d(e,s.isKanjiModeEnabled()))}}(H);var nt=i,ot=u,it=c,at=l,ut=d,st=h,ct=p,ft=v,lt=B,dt=I,ht=O,gt=P,pt=H;function vt(t,e,r){var n,o,i=t.size,a=ht.getEncodedBits(e,r);for(n=0;n<15;n++)o=1==(a>>n&1),n<6?t.set(n,8,o,!0):n<8?t.set(n+1,8,o,!0):t.set(i-15+n,8,o,!0),n<8?t.set(8,i-n-1,o,!0):n<9?t.set(8,15-n-1+1,o,!0):t.set(8,15-n-1,o,!0);t.set(i-8,8,1,!0)}function yt(t,e,r){var n=new it;r.forEach((function(e){n.put(e.mode.bit,4),n.put(e.getLength(),gt.getCharCountIndicator(e.mode,t)),e.write(n)}));var o=8*(nt.getSymbolTotalCodewords(t)-ft.getTotalCodewordsCount(t,e));for(n.getLengthInBits()+4<=o&&n.put(0,4);n.getLengthInBits()%8!=0;)n.putBit(0);for(var i=(o-n.getLengthInBits())/8,a=0;a<i;a++)n.put(a%2?17:236,8);return function(t,e,r){for(var n=nt.getSymbolTotalCodewords(e),o=ft.getTotalCodewordsCount(e,r),i=n-o,a=ft.getBlocksCount(e,r),u=a-n%a,s=Math.floor(n/a),c=Math.floor(i/a),f=c+1,l=s-c,d=new lt(l),h=0,g=new Array(a),p=new Array(a),v=0,y=new Uint8Array(t.buffer),m=0;m<a;m++){var w=m<u?c:f;g[m]=y.slice(h,h+w),p[m]=d.encode(g[m]),h+=w,v=Math.max(v,w)}var E,b,A=new Uint8Array(n),C=0;for(E=0;E<v;E++)for(b=0;b<a;b++)E<g[b].length&&(A[C++]=g[b][E]);for(E=0;E<l;E++)for(b=0;b<a;b++)A[C++]=p[b][E];return A}(n,t,e)}function mt(t,e,r,n){var o;if(Array.isArray(t))o=pt.fromArray(t);else{if("string"!=typeof t)throw new Error("Invalid data");var i=e;if(!i){var a=pt.rawSplit(t);i=dt.getBestVersionForData(a,r)}o=pt.fromString(t,i||40)}var u=dt.getBestVersionForData(o,r);if(!u)throw new Error("The amount of data is too big to be stored in a QR Code");if(e){if(e<u)throw new Error("\nThe chosen QR Code version cannot contain this amount of data.\nMinimum version required to store current data is: "+u+".\n")}else e=u;var s=yt(e,r,o),c=nt.getSymbolSize(e),f=new at(c);return function(t,e){for(var r=t.size,n=st.getPositions(e),o=0;o<n.length;o++)for(var i=n[o][0],a=n[o][1],u=-1;u<=7;u++)if(!(i+u<=-1||r<=i+u))for(var s=-1;s<=7;s++)a+s<=-1||r<=a+s||(u>=0&&u<=6&&(0===s||6===s)||s>=0&&s<=6&&(0===u||6===u)||u>=2&&u<=4&&s>=2&&s<=4?t.set(i+u,a+s,!0,!0):t.set(i+u,a+s,!1,!0))}(f,e),function(t){for(var e=t.size,r=8;r<e-8;r++){var n=r%2==0;t.set(r,6,n,!0),t.set(6,r,n,!0)}}(f),function(t,e){for(var r=ut.getPositions(e),n=0;n<r.length;n++)for(var o=r[n][0],i=r[n][1],a=-2;a<=2;a++)for(var u=-2;u<=2;u++)-2===a||2===a||-2===u||2===u||0===a&&0===u?t.set(o+a,i+u,!0,!0):t.set(o+a,i+u,!1,!0)}(f,e),vt(f,r,0),e>=7&&function(t,e){for(var r,n,o,i=t.size,a=dt.getEncodedBits(e),u=0;u<18;u++)r=Math.floor(u/3),n=u%3+i-8-3,o=1==(a>>u&1),t.set(r,n,o,!0),t.set(n,r,o,!0)}(f,e),function(t,e){for(var r=t.size,n=-1,o=r-1,i=7,a=0,u=r-1;u>0;u-=2)for(6===u&&u--;;){for(var s=0;s<2;s++)if(!t.isReserved(o,u-s)){var c=!1;a<e.length&&(c=1==(e[a]>>>i&1)),t.set(o,u-s,c),-1==--i&&(a++,i=7)}if((o+=n)<0||r<=o){o-=n,n=-n;break}}}(f,s),isNaN(n)&&(n=ct.getBestMask(f,vt.bind(null,f,r))),ct.applyMask(n,f),vt(f,r,n),{modules:f,version:e,errorCorrectionLevel:r,maskPattern:n,segments:o}}o.create=function(t,e){if(void 0===t||""===t)throw new Error("No input text");var r,n,o=ot.M;return void 0!==e&&(o=ot.from(e.errorCorrectionLevel,ot.M),r=dt.from(e.version),n=ct.from(e.maskPattern),e.toSJISFunc&&nt.setToSJISFunction(e.toSJISFunc)),mt(t,r,o,n)};var wt={},Et={};!function(t){function e(t){if("number"==typeof t&&(t=t.toString()),"string"!=typeof t)throw new Error("Color should be defined as hex string");var e=t.slice().replace("#","").split("");if(e.length<3||5===e.length||e.length>8)throw new Error("Invalid hex color: "+t);3!==e.length&&4!==e.length||(e=Array.prototype.concat.apply([],e.map((function(t){return[t,t]})))),6===e.length&&e.push("F","F");var r=parseInt(e.join(""),16);return{r:r>>24&255,g:r>>16&255,b:r>>8&255,a:255&r,hex:"#"+e.slice(0,6).join("")}}t.getOptions=function(t){t||(t={}),t.color||(t.color={});var r=void 0===t.margin||null===t.margin||t.margin<0?4:t.margin,n=t.width&&t.width>=21?t.width:void 0,o=t.scale||4;return{width:n,scale:n?4:o,margin:r,color:{dark:e(t.color.dark||"#000000ff"),light:e(t.color.light||"#ffffffff")},type:t.type,rendererOpts:t.rendererOpts||{}}},t.getScale=function(t,e){return e.width&&e.width>=t+2*e.margin?e.width/(t+2*e.margin):e.scale},t.getImageWidth=function(e,r){var n=t.getScale(e,r);return Math.floor((e+2*r.margin)*n)},t.qrToImageData=function(e,r,n){for(var o=r.modules.size,i=r.modules.data,a=t.getScale(o,n),u=Math.floor((o+2*n.margin)*a),s=n.margin*a,c=[n.color.light,n.color.dark],f=0;f<u;f++)for(var l=0;l<u;l++){var d=4*(f*u+l),h=n.color.light;f>=s&&l>=s&&f<u-s&&l<u-s&&(h=c[i[Math.floor((f-s)/a)*o+Math.floor((l-s)/a)]?1:0]),e[d++]=h.r,e[d++]=h.g,e[d++]=h.b,e[d]=h.a}}}(Et),function(t){var e=Et;t.render=function(t,r,n){var o=n,i=r;void 0!==o||r&&r.getContext||(o=r,r=void 0),r||(i=function(){try{return document.createElement("canvas")}catch(t){throw new Error("You need to specify a canvas element")}}()),o=e.getOptions(o);var a=e.getImageWidth(t.modules.size,o),u=i.getContext("2d"),s=u.createImageData(a,a);return e.qrToImageData(s.data,t,o),function(t,e,r){t.clearRect(0,0,e.width,e.height),e.style||(e.style={}),e.height=r,e.width=r,e.style.height=r+"px",e.style.width=r+"px"}(u,i,a),u.putImageData(s,0,0),i},t.renderToDataURL=function(e,r,n){var o=n;void 0!==o||r&&r.getContext||(o=r,r=void 0),o||(o={});var i=t.render(e,r,o),a=o.type||"image/png",u=o.rendererOpts||{};return i.toDataURL(a,u.quality)}}(wt);var bt={},At=Et;function Ct(t,e){var r=t.a/255,n=e+'="'+t.hex+'"';return r<1?n+" "+e+'-opacity="'+r.toFixed(2).slice(1)+'"':n}function Tt(t,e,r){var n=t+e;return void 0!==r&&(n+=" "+r),n}bt.render=function(t,e,r){var n=At.getOptions(e),o=t.modules.size,i=t.modules.data,a=o+2*n.margin,u=n.color.light.a?"<path "+Ct(n.color.light,"fill")+' d="M0 0h'+a+"v"+a+'H0z"/>':"",s="<path "+Ct(n.color.dark,"stroke")+' d="'+function(t,e,r){for(var n="",o=0,i=!1,a=0,u=0;u<t.length;u++){var s=Math.floor(u%e),c=Math.floor(u/e);s||i||(i=!0),t[u]?(a++,u>0&&s>0&&t[u-1]||(n+=i?Tt("M",s+r,.5+c+r):Tt("m",o,0),o=0,i=!1),s+1<e&&t[u+1]||(n+=Tt("h",a),a=0)):o++}return n}(i,o,n.margin)+'"/>',c='viewBox="0 0 '+a+" "+a+'"',f='<svg xmlns="http://www.w3.org/2000/svg" '+(n.width?'width="'+n.width+'" height="'+n.width+'" ':"")+c+' shape-rendering="crispEdges">'+u+s+"</svg>\n";return"function"==typeof r&&r(null,f),f};var Mt=function(){return"function"==typeof Promise&&Promise.prototype&&Promise.prototype.then},Bt=o,It=wt,Pt=bt;function Nt(t,e,r,n,o){var i=[].slice.call(arguments,1),a=i.length,u="function"==typeof i[a-1];if(!u&&!Mt())throw new Error("Callback required as last argument");if(!u){if(a<1)throw new Error("Too few arguments provided");return 1===a?(r=e,e=n=void 0):2!==a||e.getContext||(n=r,r=e,e=void 0),new Promise((function(o,i){try{var a=Bt.create(r,n);o(t(a,e,n))}catch(u){i(u)}}))}if(a<2)throw new Error("Too few arguments provided");2===a?(o=r,r=e,e=n=void 0):3===a&&(e.getContext&&void 0===o?(o=n,n=void 0):(o=n,n=r,r=e,e=void 0));try{var s=Bt.create(r,n);o(null,t(s,e,n))}catch(c){o(c)}}n.create=Bt.create,n.toCanvas=Nt.bind(null,It.render),n.toDataURL=Nt.bind(null,It.renderToDataURL),n.toString=Nt.bind(null,(function(t,e,r){return Pt.render(t,r)}));var St={exports:{}};
/*!
       * clipboard.js v2.0.11
       * https://clipboardjs.com/
       *
       * Licensed MIT © Zeno Rocha
       */!function(t,e){var r;r=function(){return function(){var t={686:function(t,e,r){r.d(e,{default:function(){return A}});var n=r(279),o=r.n(n),i=r(370),a=r.n(i),u=r(817),s=r.n(u);function c(t){try{return document.execCommand(t)}catch(e){return!1}}var f=function(t){var e=s()(t);return c("cut"),e},l=function(t,e){var r=function(t){var e="rtl"===document.documentElement.getAttribute("dir"),r=document.createElement("textarea");r.style.fontSize="12pt",r.style.border="0",r.style.padding="0",r.style.margin="0",r.style.position="absolute",r.style[e?"right":"left"]="-9999px";var n=window.pageYOffset||document.documentElement.scrollTop;return r.style.top="".concat(n,"px"),r.setAttribute("readonly",""),r.value=t,r}(t);e.container.appendChild(r);var n=s()(r);return c("copy"),r.remove(),n},d=function(t){var e=arguments.length>1&&void 0!==arguments[1]?arguments[1]:{container:document.body},r="";return"string"==typeof t?r=l(t,e):t instanceof HTMLInputElement&&!["text","search","url","tel","password"].includes(null==t?void 0:t.type)?r=l(t.value,e):(r=s()(t),c("copy")),r};function h(t){return h="function"==typeof Symbol&&"symbol"==typeof Symbol.iterator?function(t){return typeof t}:function(t){return t&&"function"==typeof Symbol&&t.constructor===Symbol&&t!==Symbol.prototype?"symbol":typeof t},h(t)}var g=function(){var t=arguments.length>0&&void 0!==arguments[0]?arguments[0]:{},e=t.action,r=void 0===e?"copy":e,n=t.container,o=t.target,i=t.text;if("copy"!==r&&"cut"!==r)throw new Error('Invalid "action" value, use either "copy" or "cut"');if(void 0!==o){if(!o||"object"!==h(o)||1!==o.nodeType)throw new Error('Invalid "target" value, use a valid Element');if("copy"===r&&o.hasAttribute("disabled"))throw new Error('Invalid "target" attribute. Please use "readonly" instead of "disabled" attribute');if("cut"===r&&(o.hasAttribute("readonly")||o.hasAttribute("disabled")))throw new Error('Invalid "target" attribute. You can\'t cut text from elements with "readonly" or "disabled" attributes')}return i?d(i,{container:n}):o?"cut"===r?f(o):d(o,{container:n}):void 0};function p(t){return p="function"==typeof Symbol&&"symbol"==typeof Symbol.iterator?function(t){return typeof t}:function(t){return t&&"function"==typeof Symbol&&t.constructor===Symbol&&t!==Symbol.prototype?"symbol":typeof t},p(t)}function v(t,e){for(var r=0;r<e.length;r++){var n=e[r];n.enumerable=n.enumerable||!1,n.configurable=!0,"value"in n&&(n.writable=!0),Object.defineProperty(t,n.key,n)}}function y(t,e){return y=Object.setPrototypeOf||function(t,e){return t.__proto__=e,t},y(t,e)}function m(t){var e=function(){if("undefined"==typeof Reflect||!Reflect.construct)return!1;if(Reflect.construct.sham)return!1;if("function"==typeof Proxy)return!0;try{return Date.prototype.toString.call(Reflect.construct(Date,[],(function(){}))),!0}catch(t){return!1}}();return function(){var r,n,o,i=w(t);if(e){var a=w(this).constructor;r=Reflect.construct(i,arguments,a)}else r=i.apply(this,arguments);return n=this,!(o=r)||"object"!==p(o)&&"function"!=typeof o?function(t){if(void 0===t)throw new ReferenceError("this hasn't been initialised - super() hasn't been called");return t}(n):o}}function w(t){return w=Object.setPrototypeOf?Object.getPrototypeOf:function(t){return t.__proto__||Object.getPrototypeOf(t)},w(t)}function E(t,e){var r="data-clipboard-".concat(t);if(e.hasAttribute(r))return e.getAttribute(r)}var b=function(t){!function(t,e){if("function"!=typeof e&&null!==e)throw new TypeError("Super expression must either be null or a function");t.prototype=Object.create(e&&e.prototype,{constructor:{value:t,writable:!0,configurable:!0}}),e&&y(t,e)}(i,t);var e,r,n,o=m(i);function i(t,e){var r;return function(t,e){if(!(t instanceof e))throw new TypeError("Cannot call a class as a function")}(this,i),(r=o.call(this)).resolveOptions(e),r.listenClick(t),r}return e=i,r=[{key:"resolveOptions",value:function(){var t=arguments.length>0&&void 0!==arguments[0]?arguments[0]:{};this.action="function"==typeof t.action?t.action:this.defaultAction,this.target="function"==typeof t.target?t.target:this.defaultTarget,this.text="function"==typeof t.text?t.text:this.defaultText,this.container="object"===p(t.container)?t.container:document.body}},{key:"listenClick",value:function(t){var e=this;this.listener=a()(t,"click",(function(t){return e.onClick(t)}))}},{key:"onClick",value:function(t){var e=t.delegateTarget||t.currentTarget,r=this.action(e)||"copy",n=g({action:r,container:this.container,target:this.target(e),text:this.text(e)});this.emit(n?"success":"error",{action:r,text:n,trigger:e,clearSelection:function(){e&&e.focus(),window.getSelection().removeAllRanges()}})}},{key:"defaultAction",value:function(t){return E("action",t)}},{key:"defaultTarget",value:function(t){var e=E("target",t);if(e)return document.querySelector(e)}},{key:"defaultText",value:function(t){return E("text",t)}},{key:"destroy",value:function(){this.listener.destroy()}}],n=[{key:"copy",value:function(t){var e=arguments.length>1&&void 0!==arguments[1]?arguments[1]:{container:document.body};return d(t,e)}},{key:"cut",value:function(t){return f(t)}},{key:"isSupported",value:function(){var t=arguments.length>0&&void 0!==arguments[0]?arguments[0]:["copy","cut"],e="string"==typeof t?[t]:t,r=!!document.queryCommandSupported;return e.forEach((function(t){r=r&&!!document.queryCommandSupported(t)})),r}}],r&&v(e.prototype,r),n&&v(e,n),i}(o()),A=b},828:function(t){if("undefined"!=typeof Element&&!Element.prototype.matches){var e=Element.prototype;e.matches=e.matchesSelector||e.mozMatchesSelector||e.msMatchesSelector||e.oMatchesSelector||e.webkitMatchesSelector}t.exports=function(t,e){for(;t&&9!==t.nodeType;){if("function"==typeof t.matches&&t.matches(e))return t;t=t.parentNode}}},438:function(t,e,r){var n=r(828);function o(t,e,r,n,o){var a=i.apply(this,arguments);return t.addEventListener(r,a,o),{destroy:function(){t.removeEventListener(r,a,o)}}}function i(t,e,r,o){return function(r){r.delegateTarget=n(r.target,e),r.delegateTarget&&o.call(t,r)}}t.exports=function(t,e,r,n,i){return"function"==typeof t.addEventListener?o.apply(null,arguments):"function"==typeof r?o.bind(null,document).apply(null,arguments):("string"==typeof t&&(t=document.querySelectorAll(t)),Array.prototype.map.call(t,(function(t){return o(t,e,r,n,i)})))}},879:function(t,e){e.node=function(t){return void 0!==t&&t instanceof HTMLElement&&1===t.nodeType},e.nodeList=function(t){var r=Object.prototype.toString.call(t);return void 0!==t&&("[object NodeList]"===r||"[object HTMLCollection]"===r)&&"length"in t&&(0===t.length||e.node(t[0]))},e.string=function(t){return"string"==typeof t||t instanceof String},e.fn=function(t){return"[object Function]"===Object.prototype.toString.call(t)}},370:function(t,e,r){var n=r(879),o=r(438);t.exports=function(t,e,r){if(!t&&!e&&!r)throw new Error("Missing required arguments");if(!n.string(e))throw new TypeError("Second argument must be a String");if(!n.fn(r))throw new TypeError("Third argument must be a Function");if(n.node(t))return function(t,e,r){return t.addEventListener(e,r),{destroy:function(){t.removeEventListener(e,r)}}}(t,e,r);if(n.nodeList(t))return function(t,e,r){return Array.prototype.forEach.call(t,(function(t){t.addEventListener(e,r)})),{destroy:function(){Array.prototype.forEach.call(t,(function(t){t.removeEventListener(e,r)}))}}}(t,e,r);if(n.string(t))return function(t,e,r){return o(document.body,t,e,r)}(t,e,r);throw new TypeError("First argument must be a String, HTMLElement, HTMLCollection, or NodeList")}},817:function(t){t.exports=function(t){var e;if("SELECT"===t.nodeName)t.focus(),e=t.value;else if("INPUT"===t.nodeName||"TEXTAREA"===t.nodeName){var r=t.hasAttribute("readonly");r||t.setAttribute("readonly",""),t.select(),t.setSelectionRange(0,t.value.length),r||t.removeAttribute("readonly"),e=t.value}else{t.hasAttribute("contenteditable")&&t.focus();var n=window.getSelection(),o=document.createRange();o.selectNodeContents(t),n.removeAllRanges(),n.addRange(o),e=n.toString()}return e}},279:function(t){function e(){}e.prototype={on:function(t,e,r){var n=this.e||(this.e={});return(n[t]||(n[t]=[])).push({fn:e,ctx:r}),this},once:function(t,e,r){var n=this;function o(){n.off(t,o),e.apply(r,arguments)}return o._=e,this.on(t,o,r)},emit:function(t){for(var e=[].slice.call(arguments,1),r=((this.e||(this.e={}))[t]||[]).slice(),n=0,o=r.length;n<o;n++)r[n].fn.apply(r[n].ctx,e);return this},off:function(t,e){var r=this.e||(this.e={}),n=r[t],o=[];if(n&&e)for(var i=0,a=n.length;i<a;i++)n[i].fn!==e&&n[i].fn._!==e&&o.push(n[i]);return o.length?r[t]=o:delete r[t],this}},t.exports=e,t.exports.TinyEmitter=e}},e={};function r(n){if(e[n])return e[n].exports;var o=e[n]={exports:{}};return t[n](o,o.exports,r),o.exports}return r.n=function(t){var e=t&&t.__esModule?function(){return t.default}:function(){return t};return r.d(e,{a:e}),e},r.d=function(t,e){for(var n in e)r.o(e,n)&&!r.o(t,n)&&Object.defineProperty(t,n,{enumerable:!0,get:e[n]})},r.o=function(t,e){return Object.prototype.hasOwnProperty.call(t,e)},r(686)}().default},t.exports=r()}(St),t("C",r(St.exports))}}}));