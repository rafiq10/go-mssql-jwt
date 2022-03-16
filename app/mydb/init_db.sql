USE [users-db] 

if not exists (select * from sysobjects where name='Users' and xtype='U')
	CREATE TABLE [dbo].[Users](
		[TF] [nchar](7) NOT NULL,
		[UserName] nvarchar(150),	
		[Salt] char(25)	,
		[Pwd] varchar(20) ,
		
	CONSTRAINT [PK_Users] PRIMARY KEY CLUSTERED 
	(
		[TF] ASC
	)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
	) ON [PRIMARY] 


	create unique index idx_Users_UserName
		on dbo.Users (UserName) include (Salt,Pwd);

if not exists (select * from sysobjects where name='Keys' and xtype='U')
	CREATE TABLE [dbo].[Keys](
		[id] INT  NOT NULL IDENTITY PRIMARY KEY,
		[key] [nvarchar](max) NOT NULL	
	)
