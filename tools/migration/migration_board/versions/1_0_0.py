"""<=1.0.0 to 1.1.0

Revision ID: 1.1.0
Revises: 

"""

# revision identifiers, used by Alembic.
revision = '1.0.0'
down_revision = None
branch_labels = None
depends_on = None

from alembic import op
#from db_meta import *
import sqlalchemy as sa

from sqlalchemy.dialects import mysql

def upgrade():
    """
    update schema&data
    """
    op.drop_column('user', 'project_admin')
    op.add_column('user', sa.Column('repo_token', mysql.VARCHAR(127), nullable=True))

def downgrade():
    """
    Downgrade has been disabled.
    """
    pass
