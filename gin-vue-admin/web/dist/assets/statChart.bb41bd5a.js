/*! 
 Build based on gin-vue-admin 
 Time : 1682153148000 */
import{a,al as e,r as t,Q as s,W as i,bt as o,F as r,o as l,c as d,d as n,a4 as c,K as x,V as b,bL as u,bh as m,bi as h}from"./index.478f69c8.js";import{a as p,b as _}from"./stat.6c08e69b.js";import{_ as f}from"./_plugin-vue_export-helper.cdc0426e.js";const v=a({data:[],data_axis:[],total:0,rank:[],rank_axis:[]}),k=async a=>{const e=await p({...a});0===e.code&&null!=e.data?(v.data=e.data.data,v.data_axis=e.data.data_axis,v.total=v.data.reduce(((a,e)=>a+e),0)):(v.data=[],v.data_axis=[],v.total=0);const t=await _({...a});0===t.code&&null!=t.data?(v.rank=t.data.rank,v.rank_axis=t.data.rank_axis):(v.rank=[],v.rank_axis=[])},y={className:"dashboard-line-box"},w={className:"dashboard-line-title"},g={className:"dashboard-line-box"},L=(a=>(m("data-v-671bc32a"),a=a(),h(),a))((()=>n("div",{className:"dashboard-line-title"}," 流量排行榜 ",-1))),N=f({__name:"statChart",setup(a){const m=e(null),h=e(null),p=t(null),_=t(null),f=v,k=a=>a>=1073741824?(a/1073741824).toFixed(1)+" G":a>=1048576?(a/1048576).toFixed(1)+" M":a>=1024?(a/1024).toFixed(1)+" K":a.toFixed(1),N=a=>{h.value.setOption({tooltip:{trigger:"axis",axisPointer:{type:"shadow"},formatter:a=>{let e=a[0].name+"<br>";return a.forEach((a=>{e+=a.marker+" "+a.seriesName+" : "+k(a.value)+"<br>"})),e}},grid:{left:"1%",right:"1%",top:"1%",bottom:"1%",containLabel:!0},xAxis:{data:a.rank_axis,axisTick:{show:!1},axisLine:{show:!1},z:10},yAxis:{axisLine:{show:!1},axisTick:{show:!1},axisLabel:{color:"#999",show:!1}},dataZoom:[{type:"inside"}],series:[{type:"bar",barWidth:"20%",itemStyle:{borderRadius:[5,5,0,0],color:"#188df0"},emphasis:{itemStyle:{color:"#188df0"}},data:a.rank}]}),m.value.setOption({tooltip:{trigger:"axis",axisPointer:{type:"shadow"},formatter:a=>{let e=a[0].name+"<br>";return a.forEach((a=>{e+=a.marker+" "+a.seriesName+" : "+k(a.value)+"<br>"})),e}},grid:{left:"1%",right:"1%",top:"1%",bottom:"1%",containLabel:!0},xAxis:{data:a.data_axis,axisTick:{show:!1},axisLine:{show:!1},z:10},yAxis:{axisLine:{show:!1},axisTick:{show:!1},axisLabel:{color:"#999",show:!1}},dataZoom:[{type:"inside"}],series:[{type:"bar",barWidth:"20%",itemStyle:{borderRadius:[5,5,0,0],color:"#188df0"},emphasis:{itemStyle:{color:"#188df0"}},data:a.data}]})};return s((async()=>{await i(),m.value=u(p.value),h.value=u(_.value),N(f)})),o((()=>{m.value&&(m.value.dispose(),m.value=null)})),r((()=>f),(a=>{N(a)}),{deep:!0}),(a,e)=>(l(),d(b,null,[n("div",y,[n("div",w," 总流量: "+c(k(x(f).total)),1),n("div",{ref_key:"echart",ref:p,className:"dashboard-line"},null,512)]),n("div",g,[L,n("div",{ref_key:"rank_echart",ref:_,className:"dashboard-line"},null,512)])],64))}},[["__scopeId","data-v-671bc32a"]]),S=Object.freeze(Object.defineProperty({__proto__:null,default:N},Symbol.toStringTag,{value:"Module"}));export{N as E,S as a,k as s};