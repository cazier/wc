import pathlib
import datetime

import yaml
import typer
import IPython
from peewee import Model, TextField, BooleanField, IntegerField, DateTimeField, SqliteDatabase, ForeignKeyField

app = typer.Typer()
conn = SqliteDatabase(None)


class BaseModel(Model):
    class Meta:
        database = conn


class Country(BaseModel):
    name = TextField()
    group = TextField()
    fifa_code = TextField()


class Match(BaseModel):
    day = IntegerField()
    played = BooleanField(default=False)

    a = ForeignKeyField(Country, backref="matches")
    b = ForeignKeyField(Country, backref="matches")
    stage = TextField()

    when = DateTimeField()
    assigned = BooleanField(default=False)


class Player(BaseModel):
    country = ForeignKeyField(Country, backref="country")

    name = TextField()
    position = TextField()
    number = IntegerField()

    goals = IntegerField()
    yellow = IntegerField()
    red = IntegerField()
    saves = IntegerField()


@app.command()
def init(
    team_yaml: pathlib.Path = typer.Argument(..., metavar="TEAMS"),
    match_yaml: pathlib.Path = typer.Argument(..., metavar="MATCHES"),
    player_yaml: pathlib.Path = typer.Argument(..., metavar="PLAYERS"),
    database: pathlib.Path = typer.Option(pathlib.Path("python.db"), help="path to the created database file"),
) -> None:
    conn.init(str(database))
    conn.create_tables([Country, Match, Player])

    storage: dict[str, Country] = {}
    teams = yaml.safe_load(team_yaml.read_text()) + [
        {"name": "Team A", "group": "", "code": "<A>"},
        {"name": "Team B", "group": "", "code": "<B>"},
    ]
    matches = yaml.safe_load(match_yaml.read_text())
    # players = yaml.safe_load(player_yaml.read_text())

    for team in teams:
        output, _ = Country.get_or_create(name=team["name"], group=team["group"], fifa_code=team["code"])
        storage[team["name"]] = output
        storage[team["code"]] = output

    for match in matches:
        date = datetime.datetime.strptime(match["date"], "%d-%b-%y")
        # date.hour, date.minute = map(int, match['time'].split(':'))
        output = Match.get_or_create(
            day=0,
            played=False,
            a=storage[match["a"]],
            b=storage[match["b"]],
            stage=match["stage"],
            when=date,
        )


@app.command()
def db(database: pathlib.Path = typer.Option(pathlib.Path("python.db"), help="path to the created database file")):
    conn.init(str(database))

    IPython.embed()


if __name__ == "__main__":
    app()
