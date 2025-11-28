# 替换为你错误的邮箱（提交时用的邮箱）
OLD_EMAIL="11730511+cheng-ruimx@user.noreply.gitee.com"
# 替换为你的 GitHub 用户名
CORRECT_NAME="mxxmstar"
# 替换为你绑定 GitHub 并验证的邮箱
CORRECT_EMAIL="2297171005@qq.com"

# 执行以下命令，批量修正历史提交
git filter-branch --env-filter '
if [ "$GIT_COMMITTER_EMAIL" = "$OLD_EMAIL" ]
then
  export GIT_COMMITTER_NAME="$CORRECT_NAME"
  export GIT_COMMITTER_EMAIL="$CORRECT_EMAIL"
fi
if [ "$GIT_AUTHOR_EMAIL" = "$OLD_EMAIL" ]
then
  export GIT_AUTHOR_NAME="$CORRECT_NAME"
  export GIT_AUTHOR_EMAIL="$CORRECT_EMAIL"
fi
' --tag-name-filter cat -- --branches --tags