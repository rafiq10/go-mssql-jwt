CREATE TABLE IF NOT EXISTS auth.keys(id text primary key, key text not null);
CREATE TABLE IF NOT EXISTS auth.cmb_departments(department text primary key);
CREATE TABLE IF NOT EXISTS auth.cmb_roles(usr_role text primary key);
CREATE TABLE IF NOT EXISTS auth.users(tf text primary key, 
					user_name text,
					email text,
					pwd text,
					created_at timestamp,
					usr_role text,
					department text,
                    session_id text
					);



DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_users_role') THEN
        ALTER TABLE auth.users
            ADD CONSTRAINT fk_users_role
            FOREIGN KEY (usr_role) REFERENCES auth.cmb_roles(usr_role);
    END IF;
END;
$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_users_department') THEN
        ALTER TABLE auth.users
            ADD CONSTRAINT fk_users_department
            FOREIGN KEY (department) REFERENCES auth.cmb_departments(department);
    END IF;
END;
$$;


insert into auth.cmb_departments(department) values ('GF'),('Comercial'),('Finanzas-Esp') on conflict  (department) do nothing;
insert into auth.cmb_roles(usr_role) values ('admin'),('tis-gf-oper'),('mgr-read') on conflict (usr_role) do nothing;
insert into auth.users(tf, user_name,email,pwd, created_at,usr_role,department,session_id) values 
('TF05079','Rafal Bil','bilrafal@gmail.com','$2a$10$nbylS5zlxR6hqGZg7cnQc.L.vi4mDFTCHzzWp1Bqn2P0mKoSsf5sG',current_timestamp,'admin','GF','')
on conflict (tf) do nothing;
insert into auth.users(tf, user_name,email,pwd, created_at,usr_role,department,session_id) values 
('TF05069','Edu MS','edums@gmail.com','$2a$10$ybMx2eDHAOUqi65lQEyBSeeBKdQHUDdvLdn64600S.5ax26VVcJKu',current_timestamp,'tis-gf-oper','GF','')
on conflict (tf) do nothing;