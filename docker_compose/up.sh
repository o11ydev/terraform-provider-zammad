#!/bin/bash

until wget -O /dev/null http://127.0.0.1:80/api/v1/getting_started
do
    sleep 1
done

/opt/zammad/bin/rails console <<< "Setting.set('system_init_done', true)"

cat <<'END' | /opt/zammad/bin/rails console
u = User.new(login: :admin, password: :admin, active: true,updated_by_id: 1, created_by_id:1, roles: Role.where(name: ['Agent', 'Admin']))
u.save!
ps = Permission.all.map do |x| x.name end
t = Token.create(user_id: u.id, persistent: true, action: 'api', label: 'Terraform', preferences: {"permission"=>ps}, created_at: Time.now, updated_at: Time.now)
t.save!
t.name = 'b9rYaoj3s2Y5dijQ3ux4TiBlexpXgYPsgEn_BiA-EQkX0o2bm1C8mDFFMqqUT8Tr'
t.save!
END
