import { useUser, useLogout } from '@/features/auth/api/auth';
import { useUsers } from '../api/users';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { LogOut, User as UserIcon, Loader2, ChevronLeft, ChevronRight, Search } from 'lucide-react';
import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useDebounce } from '@/hooks/use-debounce';

export const Dashboard = () => {
  const { user, isAuthenticated } = useUser();
  const logout = useLogout();
  const navigate = useNavigate();
  
  // State for pagination and search
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState('');
  const debouncedSearch = useDebounce(search, 500); // 500ms delay

  const limit = 10;

  // Reset page when search changes
  useEffect(() => {
    setPage(1);
  }, [debouncedSearch]);

  const { data, isLoading: usersLoading } = useUsers({
    page,
    limit,
    search: debouncedSearch,
  });

  const users = data?.data || [];
  const meta = data?.meta;
  const totalPages = meta ? Math.ceil(meta.total / limit) : 0;

  useEffect(() => {
    if (!isAuthenticated) {
      navigate('/login');
    }
  }, [isAuthenticated, navigate]);

  if (!isAuthenticated) return null;

  return (
    <div className="min-h-screen bg-gray-50 p-8">
      <div className="max-w-6xl mx-auto space-y-8">
        <div className="flex justify-between items-center">
          <div>
            <h1 className="text-3xl font-bold tracking-tight">Dashboard</h1>
            <p className="text-muted-foreground">Welcome back, {user?.name}</p>
          </div>
          <Button variant="outline" onClick={logout}>
            <LogOut className="mr-2 h-4 w-4" />
            Logout
          </Button>
        </div>

        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Your Profile</CardTitle>
              <UserIcon className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{user?.name}</div>
              <p className="text-xs text-muted-foreground">{user?.email}</p>
            </CardContent>
          </Card>
           <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Total Users</CardTitle>
              <UserIcon className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{meta?.total ? meta.total.toLocaleString() : '...'}</div>
              <p className="text-xs text-muted-foreground">Registered users</p>
            </CardContent>
          </Card>
        </div>

        <Card>
          <CardHeader>
            <div className="flex justify-between items-center">
                <div>
                    <CardTitle>All Users</CardTitle>
                    <CardDescription>List of all registered users in the system.</CardDescription>
                </div>
                <div className="relative w-64">
                    <Search className="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground" />
                    <Input 
                        placeholder="Search users..." 
                        className="pl-8" 
                        value={search}
                        onChange={(e) => setSearch(e.target.value)}
                    />
                </div>
            </div>
          </CardHeader>
          <CardContent>
            {usersLoading ? (
               <div className="flex justify-center p-8">
                 <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
               </div>
            ) : (
            <>
            <div className="rounded-md border">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>ID</TableHead>
                  <TableHead>Name</TableHead>
                  <TableHead>Email</TableHead>
                  <TableHead>Joined At</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {users.length > 0 ? users.map((u) => (
                  <TableRow key={u.id}>
                    <TableCell>{u.id}</TableCell>
                    <TableCell className="font-medium">{u.name}</TableCell>
                    <TableCell>{u.email}</TableCell>
                    <TableCell>{new Date(u.created_at).toLocaleDateString()}</TableCell>
                  </TableRow>
                )) : (
                    <TableRow>
                        <TableCell colSpan={4} className="h-24 text-center">
                            No results found.
                        </TableCell>
                    </TableRow>
                )}
              </TableBody>
            </Table>
            </div>
            
            <div className="flex items-center justify-center space-x-2 py-4">
                <Button
                    variant="outline"
                    size="sm"
                    onClick={() => setPage((p) => Math.max(1, p - 1))}
                    disabled={page === 1}
                >
                    <ChevronLeft className="h-4 w-4" />
                </Button>

                {/* Pagination Numbers */}
                {(() => {
                  const range = [];
                  const delta = 2; // How many pages to show around current page

                  // Always show first page
                  if (page > 1 + delta) {
                    range.push(1);
                    if (page > 2 + delta) range.push('...');
                  }

                  // Pages around current
                  for (let i = Math.max(1, page - delta); i <= Math.min(totalPages, page + delta); i++) {
                     range.push(i);
                  }

                  // Always show last page
                  if (page < totalPages - delta) {
                     if (page < totalPages - delta - 1) range.push('...');
                     range.push(totalPages);
                  }

                  return range.map((p, index) => (
                    typeof p === 'number' ? (
                       <Button
                          key={index}
                          variant={page === p ? "default" : "outline"}
                          size="sm"
                          onClick={() => setPage(p)}
                          className={page === p ? "pointer-events-none" : ""}
                        >
                          {p}
                        </Button>
                    ) : (
                       <span key={index} className="px-2 text-muted-foreground">...</span>
                    )
                  ));
                })()}

                <Button
                    variant="outline"
                    size="sm"
                    onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                    disabled={page === totalPages}
                >
                    <ChevronRight className="h-4 w-4" />
                </Button>
            </div>
            </>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
};
