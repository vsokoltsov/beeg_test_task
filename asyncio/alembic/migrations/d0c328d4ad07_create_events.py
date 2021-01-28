"""create events

Revision ID: d0c328d4ad07
Revises: 
Create Date: 2021-01-26 02:33:51.553080

"""
from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision = 'd0c328d4ad07'
down_revision = None
branch_labels = None
depends_on = None


def upgrade():
    query = """
    create table beeg_events (
        id int unsigned not null,
        label char(100) not null,
        count int unsigned not null default 1,
        unique key(id, label)
    );
    """
    conn = op.get_bind()
    conn.execute(query)
    


def downgrade():
    query = """
    DROP TABLE beeg_events;
    """
    conn = op.get_bind()
    conn.execute(query)
    
