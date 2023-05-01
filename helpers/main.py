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
    stage = IntegerField()

    when = DateTimeField()
    assigned = BooleanField()


class Player(BaseModel):
    country = ForeignKeyField(Countries, backref="country")

    name = TextField()
    position = TextField()
    number = IntegerField()

    goals = IntegerField()
    yellow = IntegerField()
    red = IntegerField()
    saves = IntegerField()


def main(
    team_yaml: pathlib.Path = typer.Argument(..., metavar="TEAMS"),
    match_yaml: pathlib.Path = typer.Argument(..., metavar="MATCHES"),
    player_yaml: pathlib.Path = typer.Argument(..., metavar="PLAYERS"),
    database: pathlib.Path = typer.Option(pathlib.Path("python.db"), help="path to the created database file"),
) -> None:
    db.init(str(database))
    db.create_tables([Countries, Match])

    storage = {}
    teams = yaml.safe_load(team_yaml.read_text())
    matches = yaml.safe_load(match_yaml.read_text())
    players = yaml.safe_load(player_yaml.read_text())

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

    for player in players:
        output = Player.get_or_create(name=player["name"], position=player["position"], number=player["number"])


if __name__ == "__main__":
    typer.run(main)
