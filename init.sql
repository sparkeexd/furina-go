--
-- PostgreSQL database dump
--

-- Dumped from database version 15.8
-- Dumped by pg_dump version 15.12

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: public; Type: SCHEMA; Schema: -; Owner: pg_database_owner
--

CREATE SCHEMA public;


ALTER SCHEMA public OWNER TO pg_database_owner;

--
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: pg_database_owner
--

COMMENT ON SCHEMA public IS 'standard public schema';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: hoyolab_tokens; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.hoyolab_tokens (
    user_id bigint NOT NULL,
    ltoken_v2 text NOT NULL,
    ltmid_v2 text NOT NULL,
    ltuid_v2 text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.hoyolab_tokens OWNER TO postgres;

--
-- Name: hoyolab_tokens hoyolab_tokens_pkey1; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hoyolab_tokens
    ADD CONSTRAINT hoyolab_tokens_pkey1 PRIMARY KEY (user_id);


--
-- Name: hoyolab_tokens hoyolab_tokens_user_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hoyolab_tokens
    ADD CONSTRAINT hoyolab_tokens_user_id_key UNIQUE (user_id);


--
-- Name: hoyolab_tokens; Type: ROW SECURITY; Schema: public; Owner: postgres
--

ALTER TABLE public.hoyolab_tokens ENABLE ROW LEVEL SECURITY;


--
-- PostgreSQL database dump complete
--

