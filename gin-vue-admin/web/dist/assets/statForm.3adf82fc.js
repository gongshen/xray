/*! 
 Build based on gin-vue-admin 
 Time : 1682164992000 */
import{ce as e,u as a,r as l,a as t,o,c as u,d as r,e as d,w as s,h as p,j as c,b as i,E as n}from"./index.8a0793e1.js";import{E as m,a as f}from"./form-item.de0dc8a0.js";/* empty css               *//* empty css              */import{f as b,c as v,u as y}from"./stat.161890eb.js";import"./date.23f5a973.js";import"./dictionary.a9011603.js";import"./_Uint8Array.84fc61e2.js";import"./_initCloneObject.b12e814f.js";import"./sysDictionary.5725ed3b.js";const g={class:"gva-form-box"},V={name:"Stat"},j=Object.assign(V,{setup(V){const j=e(),_=a(),w=l(""),h=l({category:"",tag:"",down:"",up:"",total:""}),U=t({category:[{required:!0,message:"",trigger:["input","blur"]}]}),k=l();(async()=>{if(j.query.id){const e=await b({ID:j.query.id});0===e.code&&(h.value=e.data.restat,w.value="update")}else w.value="create"})();const x=async()=>{var e;null==(e=k.value)||e.validate((async e=>{if(!e)return;let a;switch(w.value){case"create":default:a=await v(h.value);break;case"update":a=await y(h.value)}0===a.code&&c({type:"success",message:"创建/更改成功"})}))},q=()=>{_.go(-1)};return(e,a)=>{const l=i,t=m,c=n,b=f;return o(),u("div",null,[r("div",g,[d(b,{model:h.value,ref_key:"elFormRef",ref:k,"label-position":"right",rules:U,"label-width":"80px"},{default:s((()=>[d(t,{label:"流量分类:",prop:"category"},{default:s((()=>[d(l,{modelValue:h.value.category,"onUpdate:modelValue":a[0]||(a[0]=e=>h.value.category=e),clearable:!1,placeholder:"请输入"},null,8,["modelValue"])])),_:1}),d(t,{label:"标签:",prop:"tag"},{default:s((()=>[d(l,{modelValue:h.value.tag,"onUpdate:modelValue":a[1]||(a[1]=e=>h.value.tag=e),clearable:!0,placeholder:"请输入"},null,8,["modelValue"])])),_:1}),d(t,{label:"下行流量:",prop:"down"},{default:s((()=>[d(l,{modelValue:h.value.down,"onUpdate:modelValue":a[2]||(a[2]=e=>h.value.down=e),clearable:!0,placeholder:"请输入"},null,8,["modelValue"])])),_:1}),d(t,{label:"上行流量:",prop:"up"},{default:s((()=>[d(l,{modelValue:h.value.up,"onUpdate:modelValue":a[3]||(a[3]=e=>h.value.up=e),clearable:!0,placeholder:"请输入"},null,8,["modelValue"])])),_:1}),d(t,{label:"总流量:",prop:"total"},{default:s((()=>[d(l,{modelValue:h.value.total,"onUpdate:modelValue":a[4]||(a[4]=e=>h.value.total=e),clearable:!0,placeholder:"请输入"},null,8,["modelValue"])])),_:1}),d(t,null,{default:s((()=>[d(c,{type:"primary",onClick:x},{default:s((()=>[p("保存")])),_:1}),d(c,{type:"primary",onClick:q},{default:s((()=>[p("返回")])),_:1})])),_:1})])),_:1},8,["model","rules"])])])}}});export{j as default};