USE [users-db] 

if not exists (select * from sysobjects where name='Keys' and xtype='U')
	CREATE TABLE [dbo].[Keys](
		[id] nvarchar(50) NOT NULL PRIMARY KEY,
		[key] [nvarchar](max) NOT NULL	
	)

if not exists (select * from sysobjects where name='cmb_Departments' and xtype='U')
	CREATE TABLE [dbo].[cmb_Departments](
		[department] nvarchar(50)  NOT NULL PRIMARY KEY,
	)

if not exists (select * from sysobjects where name='cmb_Roles' and xtype='U')
	create table [dbo].[cmb_Roles](	
		[role] 				nvarchar(50) not null primary key,
	)

-- if not exists (select * from sysobjects where name='map_departmenToSiteMap' and xtype='U')
-- 	CREATE TABLE [dbo].[Departments](
-- 		[departmen] nvarchar(50),  NOT NULL PRIMARY KEY,
-- 	)

if not exists (select * from sysobjects where name='Users' and xtype='U')
	CREATE TABLE [dbo].[Users](
		[TF] 				nchar(7) 		NOT NULL PRIMARY KEY,
		[userName] 			nvarchar(150),	
		[email]				nvarchar(25),
		[salt] 				char(25)	,
		[pwd] 				nvarchar(20) ,
		[createDate] 		char(8),
		[role] 				nvarchar(50),
		[department]	  	nvarchar(50)
		
	)
	if not exists (SELECT * FROM sys.indexes WHERE name='idx_Users_UserName' AND object_id = OBJECT_ID('[dbo].Users'))
		create unique index idx_Users_UserName	on dbo.Users (UserName) include (salt,pwd)

	if not exists (SELECT * FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS WHERE CONSTRAINT_NAME ='fk_users_role')
		alter table [dbo].[Users]  with check add constraint [fk_users_role] foreign key([role])
		references [dbo].[cmb_Roles] ([role])
		alter table [dbo].[Users] check constraint [fk_users_role]
	
	if not exists (SELECT * FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS WHERE CONSTRAINT_NAME ='fk_users_department')
		alter table [dbo].[Users]  with check add constraint [fk_users_department] foreign key([department])
		references [dbo].[cmb_Departments] ([department])
		alter table [dbo].[Users] check constraint [fk_users_department]
	