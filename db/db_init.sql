CREATE DATABASE IF NOT EXISTS rede_social;
USE rede_social;

DROP TABLE IF EXISTS publicacoes;
DROP TABLE IF EXISTS seguidores;
DROP TABLE IF EXISTS usuarios;

CREATE TABLE usuarios(
    id int auto_increment primary KEY,
    nome varchar(40) not null,
    nick varchar(40) not null unique,
    email varchar(40) not null unique,
    senha varchar(100) not null,
    criadoem timestamp default current_timestamp()
) ENGINE=INNODB;

CREATE TABLE seguidores(
    usuario_id int not null,
    FOREIGN KEY (usuario_id) REFERENCES usuarios(id) ON DELETE CASCADE,
    seguidor_id int not null,
    FOREIGN KEY (seguidor_id) REFERENCES usuarios(id) ON DELETE CASCADE,
    primary key(usuario_id, seguidor_id)
) ENGINE=INNODB;

CREATE TABLE publicacoes(
    id int auto_increment primary KEY,
    titulo varchar(50) not null,
    conteudo varchar(300) not null,
    autor_id int not null,
    FOREIGN KEY (autor_id) REFERENCES usuarios(id) ON DELETE CASCADE,
    curtidas int default 0,
    criadoEm TIMESTAMP default CURRENT_TIMESTAMP
) ENGINE=INNODB;