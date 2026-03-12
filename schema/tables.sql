-- Schema definition for sqlc type generation only.
-- The actual schema is owned and managed by the dlt pipeline.

CREATE SCHEMA IF NOT EXISTS zendesk;

CREATE TABLE zendesk.tickets (
    id                              bigint PRIMARY KEY,
    subject                         varchar,
    raw_subject                     varchar,
    description                     varchar,
    status                          varchar,
    type                            varchar,
    priority                        varchar,
    recipient                       varchar,
    requester_id                    bigint,
    submitter_id                    bigint,
    assignee_id                     bigint,
    organization_id                 bigint,
    group_id                        bigint,
    brand_id                        bigint,
    ticket_form_id                  bigint,
    custom_status_id                bigint,
    external_id                     varchar,
    url                             varchar,
    created_at                      timestamptz,
    updated_at                      timestamptz,
    generated_timestamp             bigint,
    has_incidents                   boolean,
    is_public                       boolean,
    allow_channelback               boolean,
    allow_attachments               boolean,
    from_messaging_channel          boolean,
    tags                            varchar,
    custom_fields                   json,
    satisfaction_rating__score       varchar,
    satisfaction_rating__id          bigint,
    satisfaction_rating__reason      varchar,
    satisfaction_rating__reason_id   bigint,
    satisfaction_rating__comment     varchar,
    via__channel                    varchar,
    via__source__from               json,
    via__source__to                 json,
    via__source__rel                varchar,
    _dlt_load_id                    varchar,
    _dlt_id                         varchar UNIQUE NOT NULL
);

CREATE TABLE zendesk.tickets__tags (
    value           varchar,
    _dlt_root_id    varchar,
    _dlt_parent_id  varchar,
    _dlt_list_idx   bigint,
    _dlt_id         varchar UNIQUE NOT NULL
);

CREATE TABLE zendesk.tickets__collaborator_ids (
    value           bigint,
    _dlt_root_id    varchar,
    _dlt_parent_id  varchar,
    _dlt_list_idx   bigint,
    _dlt_id         varchar UNIQUE NOT NULL
);

CREATE TABLE zendesk.tickets__follower_ids (
    value           bigint,
    _dlt_root_id    varchar,
    _dlt_parent_id  varchar,
    _dlt_list_idx   bigint,
    _dlt_id         varchar UNIQUE NOT NULL
);

CREATE TABLE zendesk.tickets__email_cc_ids (
    value           bigint,
    _dlt_root_id    varchar,
    _dlt_parent_id  varchar,
    _dlt_list_idx   bigint,
    _dlt_id         varchar UNIQUE NOT NULL
);

CREATE TABLE zendesk.tickets__followup_ids (
    value           bigint,
    _dlt_root_id    varchar,
    _dlt_parent_id  varchar,
    _dlt_list_idx   bigint,
    _dlt_id         varchar UNIQUE NOT NULL
);

CREATE TABLE zendesk.users (
    id                  bigint PRIMARY KEY,
    name                varchar,
    email               varchar,
    role                varchar,
    organization_id     bigint,
    active              boolean,
    suspended           boolean,
    created_at          timestamptz,
    updated_at          timestamptz,
    time_zone           varchar,
    locale              varchar,
    phone               varchar,
    _dlt_load_id        varchar,
    _dlt_id             varchar UNIQUE NOT NULL
);

CREATE TABLE zendesk.organizations (
    id                  bigint PRIMARY KEY,
    name                varchar,
    shared_tickets      boolean,
    shared_comments     boolean,
    created_at          timestamptz,
    updated_at          timestamptz,
    details             varchar,
    notes               varchar,
    group_id            bigint,
    _dlt_load_id        varchar,
    _dlt_id             varchar UNIQUE NOT NULL
);

CREATE TABLE zendesk.groups (
    id                  bigint PRIMARY KEY,
    name                varchar,
    description         varchar,
    is_public           boolean,
    "default"           boolean,
    deleted             boolean,
    created_at          timestamptz,
    updated_at          timestamptz,
    _dlt_load_id        varchar,
    _dlt_id             varchar UNIQUE NOT NULL
);

CREATE TABLE zendesk.ticket_comments (
    id              bigint NOT NULL,
    type            varchar,
    author_id       bigint,
    body            varchar,
    html_body       varchar,
    plain_body      varchar,
    public          boolean,
    audit_id        bigint,
    created_at      timestamptz,
    ticket_id       bigint,
    via__channel    varchar,
    _dlt_load_id    varchar NOT NULL,
    _dlt_id         varchar UNIQUE NOT NULL
);
