export default {
  name: 'Name',
  script: 'Script',
  user: 'User',
  userId: 'User ID',
  userName: 'Username',
  admin: 'Admin',
  role: 'Role',
  account: 'Account',
  password: 'Password',
  mobile: 'Mobile',
  contact: 'Contact',
  project: 'Project',
  projectName: 'Project name',
  projectURL: 'Project URL',
  projectPath: 'Project path',
  autoDeploy: 'Auto deploy',
  member: 'Member',
  server: 'Server',
  serverId: 'Server ID',
  serverName: 'Server name',
  serverDescription: 'Server description',
  template: 'Template',
  package: 'Package',
  crontab: 'Crontab',
  command: 'Command',
  directory: 'Directory',
  deploy: 'Deploy',
  search: 'Search',
  tips: 'Tips',
  view: 'View',
  detail: 'Detail',
  manage: 'Manage',
  interval: 'Interval',
  desc: 'Description',
  size: 'Size',
  environment: 'Environment',
  branch: 'Branch',
  setting: 'Setting',
  baseSetting: 'Base setting',
  advancedSetting: 'Advanced setting',
  notice: 'Notice',
  task: 'Task',
  date: 'Date',
  time: 'Time',
  insertTime: 'Insert time',
  updateTime: 'Update time',
  creator: 'Creator',
  editor: 'Editor',
  op: 'Operation',
  add: 'Add',
  edit: 'Edit',
  upload: 'Upload',
  uploading: 'Uploading',
  reUpload: 'Reupload',
  import: 'Import',
  read: 'Read',
  run: 'Run',
  delete: 'Delete',
  remove: 'Remove',
  install: 'Install',
  confirm: 'Confirm',
  save: 'Save',
  cancel: 'Cancel',
  success: 'Success',
  open: 'Open',
  close: 'Close',
  stop: 'Stop',
  fail: 'Fail',
  state: 'State',
  stateOption: ['Disable', 'Enable'],
  boolOption: ['No', 'Yes'],
  runOption: ['Not run', 'Run'],
  envOption: ['Unknown', 'Production', 'Pre-release', 'Test', 'Development'],
  webhookOption: { 0: 'Nothing', 1: 'Weixin', 2: 'Dingding', 3: 'Feishu', 255: 'Custom' },
  route: {
    userProfile: 'Profile',
    deploy: 'Deploy',
    toolbox: 'Toolbox',
    json: 'JSON Formatter',
    monitor: 'APP Monitor',
    project: 'Project',
    server: 'Server',
    serverSetting: 'Server',
    template: 'Template',
    crontab: 'Crontab',
    namespace: 'Namespace',
    member: 'Member'
  },
  navbar: {
    profile: 'Your profile',
    doc: 'Document',
    logout: 'Sign out'
  },
  userPage: {
    oldPassword: 'Old password',
    newPassword: 'New password',
    rePassword: 'Confirm password'
  },
  memberPage: {
    permissionDesc: 'Admin has all namespace and project permissions',
    removeUserTips: 'This action will delete the user({name}), continue?'
  },
  namespacePage: {
    removeUserTips: `This action will delete the user's binding relationship, continue?`
  },
  templatePage: {
    scriptNotice: 'Note: The package has been uploaded to the /tmp directory',
    templateDeleteTips: 'This action will delete the template({templateName}), continue?'
  },
  serverPage: {
    sshKeyOwner: 'SSH-Key Owner',
    testConnection: 'Test',
    removeServerTips: 'This action will delete the server({serverName}), continue?'
  },
  crontabPage: {
    removeServerCrontabLabel: 'Delete crontab in the server',
    importTips: 'Please read the crontab in the server',
    selectServerTips: 'Please select the server',
    selectItemTips: 'Please select an item',
    removeCrontabServerTips: `This action will delete the server's binding relationship, continue?`
  },
  monitorPage: {
    testAppState: 'Test APP state',
    failTimes: 'failure times',
    toggleStateTips: 'This action will suspend the monitoring application({monitorName}), continue?',
    removeMontiorTips: 'This action will no longer monitor the app({monitorName}), continue?'
  },
  JSONPage: {
    expandAll: 'Expand all',
    collapseAll: 'Collapse all',
    unmarkAll: 'Unmark all',
    tips: '1.Hold down ALT and click label to achieve highlighting<br>2.Hold down SHIFT and click label to view the JSON path'
  },
  projectPage: {
    testConnection: 'Test',
    lishBranch: 'List branch',
    scriptMode: 'Script mode',
    deployNotice: 'Deploy notice',
    symlinkLabel: 'Symlink deploy',
    symlinkHeaderTips: `<p>The project synchronize to the specified directory(rsync /symlinkPath), and ln -s projectPath symlinkPath</p>
    <p>It can prevent the project from external access that are being synchronized during the process of synchronizing files</p>
    <p>Back up the latest 10 deployment files for quick rollback</p>`,
    symlinkFooterTips: `<p>If the deployment path already exists on the target server, please delete the directory manually <span style="color: red">(rm -rf projectPath)</span>, otherwise the soft chain will fail</p>
    <p>If you need to change the directory, you must manually migrate the original directory/p>`,
    afterPullScriptLabel: 'After pull script',
    afterPullScriptTips: `The script that runs on the host server after pull<br>
    For example: bash after-pull-script.sh <br>`,
    afterDeployScriptLabel: 'After deploy script',
    afterDeployScriptTips: `The script that runs on the target server after deploy<br>
    For example: bash after-deploy-script.sh<br>
    If you need to restart the service, please pay attention to whether you need nohup<br>`,
    autoDeployTitle: 'Deploy trigger: automatically build the release project after certain conditions are met',
    autoDeployTips: `Go to the webhook page in GitLab、GitHub or Gitee <br>
    Fill in URL <span style="color: red">http(s)://domian(IP)/deploy/webhook?project_id={projectId}</span><br>
    Check push event, (Gitlab can choose the corresponding branch)`,
    removeProjectTips: 'This action will delete the project({projectName}), continue?',
    removeServerTips: `This action will delete the server's({serverName}) binding relationship, continue?`,
    removeUserTips: `This action will delete the user's({userName}) binding relationship, continue?`
  },
  deployPage: {
    removeProjectTaskTips: 'This action will delete the crontab task in {projectName}, continue?',
    rollbackTips: 'This action will rebuild {commit}, continue?'
  }
}
