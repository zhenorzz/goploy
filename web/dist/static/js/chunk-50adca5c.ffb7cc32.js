(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-50adca5c"],{"24ca":function(e,t,a){"use strict";var o=a("47b7"),r=a.n(o);r.a},"47b7":function(e,t,a){},c621:function(e,t,a){"use strict";a.d(t,"c",function(){return n}),a.d(t,"d",function(){return i}),a.d(t,"a",function(){return l}),a.d(t,"b",function(){return s}),a.d(t,"e",function(){return c}),a.d(t,"f",function(){return p});var o=a("cebc"),r=a("b775");function n(e){return Object(r["a"])({url:"/template/getList",method:"get",params:Object(o["a"])({},e)})}function i(){return Object(r["a"])({url:"/template/getOption",method:"get"})}function l(e){return Object(r["a"])({url:"/template/add",method:"post",data:e})}function s(e){return Object(r["a"])({url:"/template/edit",method:"post",data:e})}function c(e){return Object(r["a"])({url:"/template/remove",method:"delete",data:{id:e}})}function p(e,t){return Object(r["a"])({url:"/template/removePackage",method:"delete",data:{templateId:e,filename:t}})}},e8da:function(e,t,a){"use strict";a.r(t);var o=function(){var e=this,t=e.$createElement,a=e._self._c||t;return a("el-row",{staticClass:"app-container"},[a("el-row",{staticClass:"app-bar",attrs:{type:"flex",justify:"end",align:"middle"}},["template"===e.activeTableName?a("el-button",{attrs:{type:"primary",icon:"el-icon-plus"},on:{click:e.handleTemplateAdd}},[e._v("添加")]):a("el-upload",{ref:"upload",attrs:{action:e.action,"before-upload":e.beforeUpload,"on-success":e.handleUploadSuccess,"on-remove":e.handleRemove,"before-remove":e.beforeRemove,"show-file-list":!1,multiple:""}},[a("el-button",{attrs:{size:"small",type:"primary"}},[e._v("点击上传")])],1)],1),e._v(" "),a("el-tabs",{staticStyle:{"box-shadow":"none"},attrs:{type:"border-card"},model:{value:e.activeTableName,callback:function(t){e.activeTableName=t},expression:"activeTableName"}},[a("el-tab-pane",{attrs:{label:"模板",name:"template"}},[a("el-table",{attrs:{border:"",stripe:"","highlight-current-row":"",data:e.templateTableData}},[a("el-table-column",{attrs:{prop:"name",label:"名称"}}),e._v(" "),a("el-table-column",{attrs:{prop:"remark",label:"描述"}}),e._v(" "),a("el-table-column",{attrs:{prop:"createTime",label:"创建时间",width:"160"}}),e._v(" "),a("el-table-column",{attrs:{prop:"updateTime",label:"更新时间",width:"160"}}),e._v(" "),a("el-table-column",{attrs:{prop:"operation",label:"操作",width:"150"},scopedSlots:e._u([{key:"default",fn:function(t){return[a("el-button",{attrs:{size:"small",type:"primary"},on:{click:function(a){return e.handleTemplateEdit(t.row)}}},[e._v("编辑")]),e._v(" "),a("el-button",{attrs:{size:"small",type:"danger"},on:{click:function(a){return e.handleTemplateDelete(t.row)}}},[e._v("删除")])]}}])})],1),e._v(" "),a("el-row",{staticStyle:{"margin-top":"10px"},attrs:{type:"flex",justify:"end"}},[a("el-pagination",{attrs:{"hide-on-single-page":"",total:e.tplPagination.total,"page-size":e.tplPagination.rows,background:"",layout:"prev, pager, next"},on:{"current-change":e.handleTplPageChange}})],1)],1),e._v(" "),a("el-tab-pane",{attrs:{label:"安装包",name:"package"}},[a("el-table",{attrs:{border:"",stripe:"","highlight-current-row":"",data:e.packageTableData}},[a("el-table-column",{attrs:{prop:"name",label:"名称"}}),e._v(" "),a("el-table-column",{attrs:{prop:"humanSize",label:"大小"}}),e._v(" "),a("el-table-column",{attrs:{prop:"createTime",label:"创建时间",width:"160"}}),e._v(" "),a("el-table-column",{attrs:{prop:"updateTime",label:"更新时间",width:"160"}}),e._v(" "),a("el-table-column",{attrs:{prop:"operation",label:"操作",width:"90"},scopedSlots:e._u([{key:"default",fn:function(t){return[a("el-upload",{ref:"upload",attrs:{action:e.action+"?packageId="+t.row.id,"before-upload":e.beforeUpload,"on-success":e.handleUploadSuccess,"on-remove":e.handleRemove,"before-remove":e.beforeRemove,"show-file-list":!1,multiple:""}},[a("el-button",{attrs:{size:"small",type:"primary"}},[e._v("重传")])],1)]}}])})],1),e._v(" "),a("el-row",{staticStyle:{"margin-top":"10px"},attrs:{type:"flex",justify:"end"}},[a("el-pagination",{attrs:{"hide-on-single-page":"",total:e.pkgPagination.total,"page-size":e.pkgPagination.rows,background:"",layout:"prev, pager, next"},on:{"current-change":e.handlePkgPageChange}})],1)],1)],1),e._v(" "),a("el-dialog",{attrs:{title:"模板设置",visible:e.dialogVisible},on:{"update:visible":function(t){e.dialogVisible=t}}},[a("el-form",{ref:"form",attrs:{rules:e.formRules,model:e.formData,"label-width":"80px"}},[a("el-form-item",{attrs:{label:"名称",prop:"name"}},[a("el-input",{attrs:{autocomplete:"off"},model:{value:e.formData.name,callback:function(t){e.$set(e.formData,"name",t)},expression:"formData.name"}})],1),e._v(" "),a("el-form-item",{attrs:{label:"描述",prop:"remark"}},[a("el-input",{attrs:{autocomplete:"off"},model:{value:e.formData.remark,callback:function(t){e.$set(e.formData,"remark",t)},expression:"formData.remark"}})],1),e._v(" "),a("el-form-item",{attrs:{label:"安装包",prop:"package"}},[a("el-select",{staticStyle:{width:"100%"},attrs:{placeholder:"选择安装包",multiple:"",clearable:"",filterable:""},model:{value:e.formData.packageIds,callback:function(t){e.$set(e.formData,"packageIds",t)},expression:"formData.packageIds"}},e._l(e.packageOption,function(e,t){return a("el-option",{key:t,attrs:{label:e.name,value:e.id}})}),1)],1),e._v(" "),a("el-form-item",{attrs:{label:"脚本",prop:"script"}},[e._v("\n        注意：安装包上传至目标服务器的/tmp目录\n        "),a("codemirror",{attrs:{options:e.cmOptions},model:{value:e.formData.script,callback:function(t){e.$set(e.formData,"script",t)},expression:"formData.script"}})],1)],1),e._v(" "),a("div",{staticClass:"dialog-footer",attrs:{slot:"footer"},slot:"footer"},[a("el-button",{on:{click:function(t){e.dialogVisible=!1}}},[e._v("取 消")]),e._v(" "),a("el-button",{attrs:{disabled:e.formProps.disabled,type:"primary"},on:{click:e.submitTemplate}},[e._v("确 定")])],1)],1)],1)},r=[],n=(a("28a5"),a("7f7f"),a("ac6a"),a("c621")),i=a("cebc"),l=a("b775");function s(e){return Object(l["a"])({url:"/package/getList",method:"get",params:Object(i["a"])({},e)})}var c=a("ed08"),p=a("8f94"),m=(a("02f0"),a("7ba3"),a("a7be"),a("1fdb"),a("4498"),{components:{codemirror:p["codemirror"]},data:function(){return{dialogVisible:!1,activeTableName:"template",templateTableData:[],tplPagination:{page:1,rows:11,total:0},packageTableData:[],pkgPagination:{page:1,rows:11,total:0},packageOption:[],action:"/package/upload",tempFormData:{},cmOptions:{tabSize:4,mode:"text/x-sh",lineNumbers:!0,line:!0,scrollbarStyle:"overlay",theme:"darcula"},formProps:{disabled:!1},formData:{id:0,name:"",remark:"",packageIds:[],packageIdStr:"",script:""},formRules:{name:[{required:!0,message:"名称",trigger:"blur"}],script:[{required:!0,message:"请输入脚本",trigger:"blur"}]}}},created:function(){this.storeFormData(),this.getTemplateList(),this.getPackageList()},methods:{getTemplateList:function(){var e=this;Object(n["c"])(this.tplPagination).then(function(t){var a=t.data.templateList||[];a.forEach(function(e){e.createTime=Object(c["b"])(e.createTime),e.updateTime=Object(c["b"])(e.updateTime)}),e.templateTableData=a,e.tplPagination=t.data.pagination})},handleTplPageChange:function(e){this.tplPagination.page=e,this.getTemplateList()},getPackageList:function(){var e=this;s(this.pkgPagination).then(function(t){var a=t.data.packageList||[];a.forEach(function(e){e.createTime=Object(c["b"])(e.createTime),e.updateTime=Object(c["b"])(e.updateTime),e.humanSize=Object(c["a"])(e.size)}),e.packageOption=e.packageTableData=a,e.pkgPagination=t.data.pagination})},handlePkgPageChange:function(e){this.pkgPagination.page=e,this.getPackageList()},beforeUpload:function(e){this.$message.info("正在上传")},handleUploadSuccess:function(e,t,a){0!==e.code?this.$message.error(e.message):(this.$message.success("上传成功"),this.getPackageList())},handleRemove:function(e,t){console.log(e,t)},beforeRemove:function(e,t){return Object(n["f"])(this.formData.id,e.name).then(function(e){return Promise.resolve(e)}).catch(function(e){return Promise.reject(e)})},handleTemplateAdd:function(){this.formProps.fileList=[],this.restoreFormData(),this.dialogVisible=!0},handleTemplateEdit:function(e){this.formData=Object.assign(this.formData,e),this.formData.packageIds=e.packageIdStr.split(",").map(function(e){return parseInt(e)}),this.dialogVisible=!0},handleTemplateDelete:function(e){var t=this;this.$confirm("此操作将删除该模板, 是否继续?","提示",{confirmButtonText:"确定",cancelButtonText:"取消",type:"warning"}).then(function(){Object(n["e"])(e.id).then(function(e){t.$message({message:e.message,type:"success",duration:5e3}),t.getList()})}).catch(function(){t.$message({type:"info",message:"已取消删除"})})},submitTemplate:function(){var e=this;this.$refs.form.validate(function(t){if(!t)return!1;e.formData.packageIdStr=e.formData.packageIds.join(","),0===e.formData.id?e.addTemplate():e.editTemplate()})},addTemplate:function(){var e=this;this.formProps.disabled=!0,Object(n["a"])(this.formData).then(function(t){e.getTemplateList(),e.$message({message:t.message,type:"success",duration:5e3})}).finally(function(){e.formProps.disabled=e.dialogVisible=!1})},editTemplate:function(){var e=this;this.formProps.disabled=!0,Object(n["b"])(this.formData).then(function(t){e.getTemplateList(),e.$message({message:t.message,type:"success",duration:5e3})}).finally(function(){e.formProps.disabled=e.dialogVisible=!1})},storeFormData:function(){this.tempFormData=JSON.parse(JSON.stringify(this.formData))},restoreFormData:function(){this.formData=JSON.parse(JSON.stringify(this.tempFormData))}}}),u=m,d=(a("24ca"),a("2877")),f=Object(d["a"])(u,o,r,!1,null,null,null);t["default"]=f.exports},ed08:function(e,t,a){"use strict";a.d(t,"b",function(){return r}),a.d(t,"a",function(){return n});a("ac6a"),a("c5f6"),a("28a5"),a("a481"),a("6b54");var o=a("7618");function r(e,t){if(0===arguments.length)return null;var a,r=t||"{y}-{m}-{d} {h}:{i}:{s}";"object"===Object(o["a"])(e)?a=e:("string"===typeof e&&/^[0-9]+$/.test(e)&&(e=parseInt(e)),"number"===typeof e&&10===e.toString().length&&(e*=1e3),a=new Date(e));var n={y:a.getFullYear(),m:a.getMonth()+1,d:a.getDate(),h:a.getHours(),i:a.getMinutes(),s:a.getSeconds(),a:a.getDay()},i=r.replace(/{(y|m|d|h|i|s|a)+}/g,function(e,t){var a=n[t];return"a"===t?["日","一","二","三","四","五","六"][a]:(e.length>0&&a<10&&(a="0"+a),a||0)});return i}function n(e){if(0===e)return"0 B";var t=1024,a=["B","KB","MB","GB","TB","PB","EB","ZB","YB"],o=Math.floor(Math.log(e)/Math.log(t));return(e/Math.pow(t,o)).toPrecision(3)+" "+a[o]}}}]);