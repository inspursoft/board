import sqlalchemy as sa
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker, relationship
from sqlalchemy.dialects import mysql

Base = declarative_base()

class User(Base):
    __tablename__ = 'user'

    id = sa.Column(sa.Integer, primary_key=True)
    username = sa.Column(sa.String(255), unique=True)
    password = sa.Column(sa.String(40), nullable=False)
    email = sa.Column(sa.String(255), unique=True)
    realname = sa.Column(sa.String(255), nullable=False)
    comment = sa.Column(sa.String(30))
    deleted = sa.Column(sa.Integer, nullable=False, server_default=sa.text("'0'"))
    system_admin = sa.Column(sa.Integer, nullable=False, server_default=sa.text("'0'"))
    reset_uuid = sa.Column(sa.String(255))
    salt = sa.Column(sa.String(255))
    repo_token = sa.Column(sa.String(127))
    creation_time = sa.Column(mysql.DATETIME)
    update_time = sa.Column(mysql.DATETIME)