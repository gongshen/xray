/*! 
 Build based on gin-vue-admin 
 Time : 1682153148000 */
import{ce as e,u as a,bf as o,b$ as t,r,F as s,bt as l,cf as c,o as n,c as i,e as u,w as d,aE as f,V as p,aD as m,K as v,a1 as b,f as h,a3 as x}from"./index.478f69c8.js";import{E as g}from"./scrollbar.f3b871cb.js";/* empty css                *//* empty css               */import k from"./index.7c7122a4.js";import{b as j}from"./index.b8ee75bc.js";import"./menuItem.4221fad5.js";/* empty css             */import"./_plugin-vue_export-helper.cdc0426e.js";import"./index.180d1878.js";import"./index.89c5aafe.js";import"./asyncSubmenu.31ecaf3d.js";const y={name:"Aside"},T=Object.assign(y,{setup(y){const T=e(),M=a(),_=o(),w=t(),B=r({}),q=()=>{switch(_.sideMode){case"#fff":B.value={background:"#fff",activeBackground:"var(--el-color-primary)",activeText:"#fff",normalText:"#333",hoverBackground:"rgba(64, 158, 255, 0.08)",hoverText:"#333"};break;case"#191a23":B.value={background:"#191a23",activeBackground:"var(--el-color-primary)",activeText:"#fff",normalText:"#fff",hoverBackground:"rgba(64, 158, 255, 0.08)",hoverText:"#fff"}}};q();const E=r("");s((()=>T),(()=>{E.value=T.meta.activeName||T.name}),{deep:!0}),s((()=>_.sideMode),(()=>{q()}));const O=r(!1);(()=>{E.value=T.meta.activeName||T.name;document.body.clientWidth<1e3&&(O.value=!O.value),c.on("collapse",(e=>{O.value=e}))})(),l((()=>{c.off("collapse")}));const N=(e,a,o,t)=>{var r,s;const l={},c={};(null==(r=w.routeMap[e])?void 0:r.parameters)&&(null==(s=w.routeMap[e])||s.parameters.forEach((e=>{"query"===e.type?l[e.key]=e.value:c[e.key]=e.value}))),e!==T.name&&(e.indexOf("http://")>-1||e.indexOf("https://")>-1?window.open(e):M.push({name:e,query:l,params:c}))};return(e,a)=>{const o=j,t=g;return n(),i("div",{style:x({background:v(_).sideMode})},[u(t,{style:{height:"calc(100vh - 60px)"}},{default:d((()=>[u(f,{duration:{enter:800,leave:100},mode:"out-in",name:"el-fade-in-linear"},{default:d((()=>[u(o,{collapse:O.value,"collapse-transition":!1,"default-active":E.value,"background-color":B.value.background,"active-text-color":B.value.active,class:"el-menu-vertical","unique-opened":"",onSelect:N},{default:d((()=>[(n(!0),i(p,null,m(v(w).asyncRouters[0].children,(e=>(n(),i(p,null,[e.hidden?h("",!0):(n(),b(k,{key:e.name,"is-collapse":O.value,"router-info":e,theme:B.value},null,8,["is-collapse","router-info","theme"]))],64)))),256))])),_:1},8,["collapse","default-active","background-color","active-text-color"])])),_:1})])),_:1})],4)}}});export{T as default};