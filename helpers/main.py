import pathlib
import datetime

import yaml
import typer
from peewee import Model, TextField, BooleanField, IntegerField, DateTimeField, SqliteDatabase, ForeignKeyField

db = SqliteDatabase(None)


class BaseModel(Model):
    class Meta:
        database = db


class Countries(BaseModel):
    name = TextField()
    group = TextField()
    fifa_code = TextField()


class Match(BaseModel):
    day = IntegerField()
    played = BooleanField()

    a = ForeignKeyField(Countries, backref="matches")
    b = ForeignKeyField(Countries, backref="matches")

    group = TextField()
    when = DateTimeField()


def main(
    team_yaml: pathlib.Path = typer.Argument(..., metavar="TEAMS"),
    match_yaml: pathlib.Path = typer.Argument(..., metavar="MATCHES"),
    database: pathlib.Path = typer.Option(pathlib.Path("python.db"), help="path to the created database file"),
) -> None:
    db.init(str(database))
    db.create_tables([Countries, Match])

    storage = {}
    teams = yaml.safe_load(team_yaml.read_text())
    matches = yaml.safe_load(match_yaml.read_text())

    for team in teams:
        output, _ = Countries.get_or_create(name=team["name"], group=team["group"], fifa_code=team["code"])
        storage[team["name"]] = output

    for match in matches:
        date = datetime.datetime.strptime(match["date"], "%d-%b-%y")
        output = Match.get_or_create(
            day=0,
            played=False,
            a=storage[match["a"]],
            b=storage[match["b"]],
            group=match["group"],
            when=date,
        )


if __name__ == "__main__":
    typer.run(main)
