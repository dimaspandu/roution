from dataclasses import dataclass, field

from roution.pattern import parse_pattern, SEGMENT_STATIC, SEGMENT_PARAM


@dataclass
class StaticRoute:
    route: str
    value: object


@dataclass
class DynamicRoute:
    route: str
    value: object
    segments: list
    param_names: list


@dataclass
class CompiledRoutes:
    static_routes: dict = field(default_factory=dict)
    dynamic_routes: list = field(default_factory=list)
    wildcard: object = None


def compile(routes):
    """Compile route definitions into an internal, reusable representation.

    Routes are partitioned into static routes (exact lookup), dynamic routes
    (segments containing parameters), and an optional wildcard fallback.
    Compilation happens once; the result is reused for every match.
    """
    static_routes = {}
    dynamic_routes = []
    wildcard = None

    for pattern, value in routes.items():
        parsed = parse_pattern(pattern)

        if parsed.is_wildcard:
            wildcard = StaticRoute(route=pattern, value=value)
            continue

        only_static = all(segment.type == SEGMENT_STATIC for segment in parsed.segments)
        if only_static:
            static_routes[parsed.raw] = StaticRoute(route=pattern, value=value)
            continue

        param_names = [
            segment.name for segment in parsed.segments if segment.type == SEGMENT_PARAM
        ]
        dynamic_routes.append(
            DynamicRoute(
                route=pattern,
                value=value,
                segments=parsed.segments,
                param_names=param_names,
            )
        )

    return CompiledRoutes(
        static_routes=static_routes,
        dynamic_routes=dynamic_routes,
        wildcard=wildcard,
    )
