/*! 
 Build based on gin-vue-admin 
 Time : 1682164992000 */
import{ce as e,u as a,r as l,a as s,o as r,c as t,d as i,e as o,w as u,h as d,j as n,b as c,E as p}from"./index.8a0793e1.js";import{E as m,a as f}from"./form-item.de0dc8a0.js";/* empty css               *//* empty css              */import{f as v,c as b,u as _}from"./binding.f4e8fb2c.js";import"./date.23f5a973.js";import"./dictionary.a9011603.js";import"./_Uint8Array.84fc61e2.js";import"./_initCloneObject.b12e814f.js";import"./sysDictionary.5725ed3b.js";const y={class:"gva-form-box"},j={name:"Binding"},g=Object.assign(j,{setup(j){const g=e(),k=a(),h=l(""),w=l({server_id:"",user_id:""}),V=s({}),x=l();(async()=>{if(g.query.id){const e=await v({ID:g.query.id});0===e.code&&(w.value=e.data.rebinding,h.value="update")}else h.value="create"})();const C=async()=>{var e;null==(e=x.value)||e.validate((async e=>{if(!e)return;let a;switch(h.value){case"create":default:a=await b(w.value);break;case"update":a=await _(w.value)}0===a.code&&n({type:"success",message:"创建/更改成功"})}))},U=()=>{k.go(-1)};return(e,a)=>{const l=c,s=m,n=p,v=f;return r(),t("div",null,[i("div",y,[o(v,{model:w.value,ref_key:"elFormRef",ref:x,"label-position":"right",rules:V,"label-width":"80px"},{default:u((()=>[o(s,{label:"服务器ip:",prop:"server_ip"},{default:u((()=>[o(l,{modelValue:w.value.server_ip,"onUpdate:modelValue":a[0]||(a[0]=e=>w.value.server_ip=e),clearable:!0,placeholder:"请输入"},null,8,["modelValue"])])),_:1}),o(s,{label:"用户名:",prop:"nick_name"},{default:u((()=>[o(l,{modelValue:w.value.nick_name,"onUpdate:modelValue":a[1]||(a[1]=e=>w.value.nick_name=e),clearable:!0,placeholder:"请输入"},null,8,["modelValue"])])),_:1}),o(s,null,{default:u((()=>[o(n,{type:"primary",onClick:C},{default:u((()=>[d("保存")])),_:1}),o(n,{type:"primary",onClick:U},{default:u((()=>[d("返回")])),_:1})])),_:1})])),_:1},8,["model","rules"])])])}}});export{g as default};