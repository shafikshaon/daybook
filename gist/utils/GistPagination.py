from collections import OrderedDict

from rest_framework.pagination import LimitOffsetPagination, _positive_int
from rest_framework.response import Response


class GistPagination(LimitOffsetPagination):
    def get_paginated_response(self, data):
        return Response(OrderedDict([
            ('count', self.count),
            ('next', self.get_next_link()),
            ('previous', self.get_previous_link()),
            ('remaining_count',
             self.count - (self.offset + self.limit) if self.offset + self.limit < self.count else 0),
            ('next_offset', self.offset + self.limit if self.offset + self.limit < self.count else 0),
            ('current_page', int(self.offset / self.limit) + 1),
            ('total_page', int(self.count / self.limit) + (0 if (self.count % self.limit == 0) else 1)),
            ('results', data)
        ]))

    def get_limit(self, request):
        if request.GET.get('disable_pagination', False):
            return 9999999  # self.max_limit
        if self.limit_query_param:
            try:
                return _positive_int(
                    request.query_params[self.limit_query_param],
                    cutoff=self.max_limit
                )
            except (KeyError, ValueError):
                pass

        return self.default_limit
