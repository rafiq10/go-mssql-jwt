insert into auth.cmb_departments(department) values ('Test-Dep') on conflict  (department) do nothing;
-- delete from auth.cmb_departments where department='Test-Dep';