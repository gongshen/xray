/*! 
 Build based on gin-vue-admin 
 Time : 1682164992000 */
import{q as e,t,B as a,D as s,A as r,bM as n,bN as o,az as l,bO as i,b2 as u,M as c,y as p,o as d,c as f,n as h,K as y,d as g,a3 as v,J as k,a4 as b,f as m,a1 as x,w as $,a2 as w,aw as N,_ as D,a6 as I}from"./index.8a0793e1.js";const B=e({type:{type:String,default:"line",values:["line","circle","dashboard"]},percentage:{type:Number,default:0,validator:e=>e>=0&&e<=100},status:{type:String,default:"",values:["","success","exception","warning"]},indeterminate:{type:Boolean,default:!1},duration:{type:Number,default:3},strokeWidth:{type:Number,default:6},strokeLinecap:{type:t(String),default:"round"},textInside:{type:Boolean,default:!1},width:{type:Number,default:126},showText:{type:Boolean,default:!0},color:{type:t([String,Array,Function]),default:""},format:{type:t(Function),default:e=>`${e}%`}}),S=["aria-valuenow"],T={viewBox:"0 0 100 100"},F=["d","stroke","stroke-width"],M=["d","stroke","opacity","stroke-linecap","stroke-width"],W={key:0},_=a({name:"ElProgress"});const z=I(D(a({..._,props:B,setup(e){const t=e,a={success:"#13ce66",exception:"#ff4949",warning:"#e6a23c",default:"#20a0ff"},D=s("progress"),I=r((()=>({width:`${t.percentage}%`,animationDuration:`${t.duration}s`,backgroundColor:O(t.percentage)}))),B=r((()=>(t.strokeWidth/t.width*100).toFixed(1))),_=r((()=>["circle","dashboard"].includes(t.type)?Number.parseInt(""+(50-Number.parseFloat(B.value)/2),10):0)),z=r((()=>{const e=_.value,a="dashboard"===t.type;return`\n          M 50 50\n          m 0 ${a?"":"-"}${e}\n          a ${e} ${e} 0 1 1 0 ${a?"-":""}${2*e}\n          a ${e} ${e} 0 1 1 0 ${a?"":"-"}${2*e}\n          `})),A=r((()=>2*Math.PI*_.value)),E=r((()=>"dashboard"===t.type?.75:1)),L=r((()=>`${-1*A.value*(1-E.value)/2}px`)),P=r((()=>({strokeDasharray:`${A.value*E.value}px, ${A.value}px`,strokeDashoffset:L.value}))),j=r((()=>({strokeDasharray:`${A.value*E.value*(t.percentage/100)}px, ${A.value}px`,strokeDashoffset:L.value,transition:"stroke-dasharray 0.6s ease 0s, stroke 0.6s ease, opacity ease 0.6s"}))),q=r((()=>{let e;return e=t.color?O(t.percentage):a[t.status]||a.default,e})),C=r((()=>"warning"===t.status?n:"line"===t.type?"success"===t.status?o:l:"success"===t.status?i:u)),J=r((()=>"line"===t.type?12+.4*t.strokeWidth:.111111*t.width+2)),K=r((()=>t.format(t.percentage)));const O=e=>{var a;const{color:s}=t;if(c(s))return s(e);if(p(s))return s;{const t=function(e){const t=100/e.length;return e.map(((e,a)=>p(e)?{color:e,percentage:(a+1)*t}:e)).sort(((e,t)=>e.percentage-t.percentage))}(s);for(const a of t)if(a.percentage>e)return a.color;return null==(a=t[t.length-1])?void 0:a.color}};return(e,t)=>(d(),f("div",{class:h([y(D).b(),y(D).m(e.type),y(D).is(e.status),{[y(D).m("without-text")]:!e.showText,[y(D).m("text-inside")]:e.textInside}]),role:"progressbar","aria-valuenow":e.percentage,"aria-valuemin":"0","aria-valuemax":"100"},["line"===e.type?(d(),f("div",{key:0,class:h(y(D).b("bar"))},[g("div",{class:h(y(D).be("bar","outer")),style:v({height:`${e.strokeWidth}px`})},[g("div",{class:h([y(D).be("bar","inner"),{[y(D).bem("bar","inner","indeterminate")]:e.indeterminate}]),style:v(y(I))},[(e.showText||e.$slots.default)&&e.textInside?(d(),f("div",{key:0,class:h(y(D).be("bar","innerText"))},[k(e.$slots,"default",{percentage:e.percentage},(()=>[g("span",null,b(y(K)),1)]))],2)):m("v-if",!0)],6)],6)],2)):(d(),f("div",{key:1,class:h(y(D).b("circle")),style:v({height:`${e.width}px`,width:`${e.width}px`})},[(d(),f("svg",T,[g("path",{class:h(y(D).be("circle","track")),d:y(z),stroke:`var(${y(D).cssVarName("fill-color-light")}, #e5e9f2)`,"stroke-width":y(B),fill:"none",style:v(y(P))},null,14,F),g("path",{class:h(y(D).be("circle","path")),d:y(z),stroke:y(q),fill:"none",opacity:e.percentage?1:0,"stroke-linecap":e.strokeLinecap,"stroke-width":y(B),style:v(y(j))},null,14,M)]))],6)),!e.showText&&!e.$slots.default||e.textInside?m("v-if",!0):(d(),f("div",{key:2,class:h(y(D).e("text")),style:v({fontSize:`${y(J)}px`})},[k(e.$slots,"default",{percentage:e.percentage},(()=>[e.status?(d(),x(y(N),{key:1},{default:$((()=>[(d(),x(w(y(C))))])),_:1})):(d(),f("span",W,b(y(K)),1))]))],6))],10,S))}}),[["__file","/home/runner/work/element-plus/element-plus/packages/components/progress/src/progress.vue"]]));export{z as E};