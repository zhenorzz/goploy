export default {
  name: '名称',
  script: '脚本',
  user: '用户',
  userId: '用户ID',
  userName: '用户名称',
  admin: '超管',
  role: '角色',
  account: '账号',
  password: '密码',
  mobile: '手机号码',
  contact: '联系方式',
  project: '项目',
  projectName: '项目名称',
  projectURL: '项目地址',
  projectPath: '部署路径',
  autoDeploy: '自动部署',
  member: '成员',
  server: '服务器',
  serverId: '服务器ID',
  serverName: '服务器名称',
  serverDescription: '服务器描述',
  template: '模板',
  package: '安装包',
  crontab: '定时',
  command: '命令',
  directory: '目录',
  deploy: '构建',
  search: '搜索',
  tips: '提示',
  view: '查看',
  detail: '详情',
  manage: '管理',
  interval: '间隔',
  desc: '描述',
  size: '大小',
  environment: '环境',
  branch: '分支',
  setting: '设置',
  baseSetting: '基本配置',
  advancedSetting: '高级配置',
  notice: '通知',
  task: '任务',
  date: '日期',
  time: '时间',
  insertTime: '创建时间',
  updateTime: '更新时间',
  creator: '创建人',
  editor: '修改人',
  op: '操作',
  add: '添加',
  edit: '编辑',
  upload: '上传',
  uploading: '上传中',
  reUpload: '重传',
  import: '导入',
  read: '读取',
  run: '运行',
  delete: '删除',
  remove: '移除',
  install: '安装',
  confirm: '确认',
  save: '保存',
  cancel: '取消',
  success: '成功',
  open: '开启',
  close: '关闭',
  stop: '暂停',
  fail: '失败',
  state: '状态',
  stateOption: ['失效', '生效'],
  boolOption: ['否', '是'],
  runOption: ['未运行', '已运行'],
  envOption: ['未知', '生产环境', '预发布环境', '测试环境', '开发环境'],
  webhookOption: { 0: '无', 1: '企业微信', 2: '钉钉', 3: '飞书', 255: '自定义' },
  route: {
    userProfile: '个人信息',
    deploy: '构建发布',
    toolbox: '工具箱',
    json: 'JSON格式化',
    monitor: '应用监控',
    project: '项目设置',
    server: '服务器管理',
    serverSetting: '服务器设置',
    template: '模板设置',
    crontab: 'Crontab管理',
    namespace: '空间设置',
    member: '成员设置'
  },
  navbar: {
    profile: '个人中心',
    doc: '帮助文档',
    logout: '退出'
  },
  userPage: {
    oldPassword: '原密码',
    newPassword: '新密码',
    rePassword: '确认新密码'
  },
  memberPage: {
    permissionDesc: '超管具有所有空间和项目权限',
    removeUserTips: '此操作将删除{name}, 是否继续?'
  },
  namespacePage: {
    removeUserTips: '此操作将永久删除该用户的绑定关系, 是否继续?'
  },
  templatePage: {
    scriptNotice: '注意：安装包上传至目标服务器的/tmp目录',
    templateDeleteTips: '此操作将删除模板({templateName}), 是否继续?'
  },
  serverPage: {
    sshKeyOwner: 'SSH-Key 所有者',
    testConnection: '测试连接',
    removeServerTips: '此操作将删除服务器({serverName}), 是否继续?'
  },
  crontabPage: {
    removeServerCrontabLabel: '删除服务器Crontab任务',
    importTips: '请读取服务器Crontab任务',
    selectServerTips: '请先选择服务器',
    selectItemTips: '请先选择需要导入的条目',
    removeCrontabServerTips: '此操作将永久删除该服务器的绑定关系, 是否继续?'
  },
  monitorPage: {
    testAppState: '测试应用状态',
    failTimes: '连续失败次数',
    toggleStateTips: '此操作将暂停监控应用({monitorName}), 是否继续?',
    removeMontiorTips: '此操作将不再监控应用({monitorName}), 是否继续?'
  },
  JSONPage: {
    expandAll: '展开所有',
    collapseAll: '收起所有',
    unmarkAll: '取消高亮',
    tips: '1.按住ALT点击label可以实现高亮<br>2.按住SHIFT可以查看JSON路径'
  },
  projectPage: {
    testConnection: '测试连接',
    lishBranch: '列出分支',
    scriptMode: '脚本类型',
    deployNotice: '构建通知',
    symlinkLabel: '软链部署(推荐)',
    symlinkHeaderTips: `<p>项目先同步到指定目录(rsync 软链目录)，然后ln -s 部署路径 软链目录</p>
    <p>可以避免项目在同步传输文件的过程中，外部访问到部分正在同步的文件</p>
    <p>备份最近10次的部署文件，以便快速回滚</p>`,
    symlinkFooterTips: `<p>如果部署路径已存在在目标服务器，请手动删除该目录<span style="color: red">rm -rf 部署路径</span>，否则软链将会不成功</p>
    <p>如需更换目录，务必手动迁移原先的目录到目标服务器</p>`,
    afterPullScriptLabel: '拉取后运行脚本',
    afterPullScriptTips: `拉取代码后在宿主服务器运行的脚本<br>
    运行方式：打包成一份脚本文件<br>
    检查服务器是否安装该脚本类型(默认以bash运行)<br>`,
    afterDeployScriptLabel: '部署后运行脚本',
    afterDeployScriptTips: `部署后在目标服务器运行的脚本<br>
    运行方式：打包成一份脚本文件<br>
    如需重启服务，请注意是否需要nohup<br>
    检查服务器是否安装该脚本类型(默认以bash运行)`,
    autoDeployTitle: '构建触发器：达成某种条件后自动构建发布项目',
    autoDeployTips: `前往GitLab、GitHub或Gitee的webhook（可前往谷歌查找各自webhook所在的位置）<br>
    填入连接<span style="color: red">http(s)://域名(IP)/deploy/webhook?project_id={projectId}</span><br>
    勾选push event即可, (Gitlab可以选对应的分支)`,
    removeProjectTips: '此操作将删除项目({projectName}), 是否继续?',
    removeServerTips: '此操作将永久删除服务器({serverName})的绑定关系, 是否继续?',
    removeUserTips: '此操作将永久删除用户({userName})的绑定关系, 是否继续?'
  },
  deployPage: {
    removeProjectTaskTips: '此操作删除{projectName}的定时任务, 是否继续?',
    rollbackTips: '此操作将重新构建{commit}, 是否继续?'
  }
}
