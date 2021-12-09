# Usage problem

# The script cannot restart the application service after deployment

- Check if the script is correct to kill the application service
- Restart the application using nohup

# $HOME is not defined

- Check if exists Env var $HOME
- If use systemd or supervisord，put the $HOME to the configure file

# The project script appear xx is not defined

- If the command was installed after the goploy started, need to restart goploy
- If not, try the absolute path or source /etc/profile